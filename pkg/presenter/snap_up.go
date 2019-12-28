package presenter

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/labstack/echo"
	add "github.com/webronin26/snap-up/pkg/usecases/add_record"
	query "github.com/webronin26/snap-up/pkg/usecases/query_product"
	update "github.com/webronin26/snap-up/pkg/usecases/update_product"
)

var (
	ProductID     int32     // 目前 產品的 ID
	MaxBuy        int32     // 一次能購買的最大量
	ItemNumber    int32     // 目前的產品剩餘數
	AvailableTime time.Time // 開賣時間

	MaxGoroutineNumber     int // 最大 goroutine 的數字
	CurrentGoroutineNumber int // 目前 goroutine 的數字

	MissionChannel chan int32 // 任務隊列
	ResultChannel  chan int32 // 結果隊列
)

type Input struct {
	CustomerID int32
	Number     int32
}

// 初始化產品
func InitSnapUp(productID int32) {

	product, err := query.Exec(productID)
	if err != nil {
		panic(fmt.Errorf("init snap-up error %s \n", err))
	}

	ProductID = productID
	MaxBuy = product.MaxBuy
	ItemNumber = product.ItemNumber
	AvailableTime = product.AvailableTime

	// 最多只開啟跟目前 CPU 數量相等的 goroutine
	MaxGoroutineNumber = runtime.NumCPU()
	CurrentGoroutineNumber = 0

	MissionChannel = make(chan int32, MaxGoroutineNumber)
	ResultChannel = make(chan int32, MaxGoroutineNumber)

	go runLooper()
}

// 啟動 Looper
func runLooper() {

	for {
		for {
			// 如果結果列隊是空的，就進到下一步驟
			if len(ResultChannel) == 0 {
				break
			}
			resultNumber := <-ResultChannel
			// 現有產品數量加上從 channel 拿出來的數字
			ItemNumber = ItemNumber + resultNumber
			CurrentGoroutineNumber = CurrentGoroutineNumber - 1
		}

		// 將產品的資料庫數量更新一下
		// 但是這個更新的數量不一定是正確的剩餘數量，很可能有其他線程把一些資源先扣掉了
		var updateInput update.Input
		updateInput.ProductID = ProductID
		updateInput.ItemNumber = ItemNumber

		if err := update.Exec(updateInput); err != nil {
			// 這邊可以輸出日誌
			fmt.Println("updating error %s \n", err)
		}

		for {
			if CurrentGoroutineNumber == MaxGoroutineNumber {
				break
			}

			if CurrentGoroutineNumber == 0 && ItemNumber == 0 {
				break
			}

			if CurrentGoroutineNumber == 0 && ItemNumber > 0 && ItemNumber < MaxBuy {
				MissionChannel <- ItemNumber
				ItemNumber = 0
				CurrentGoroutineNumber = CurrentGoroutineNumber + 1
				break
			}

			if CurrentGoroutineNumber == 0 && ItemNumber > 0 && ItemNumber >= MaxBuy {
				MissionChannel <- MaxBuy
				ItemNumber = ItemNumber - MaxBuy
				CurrentGoroutineNumber = CurrentGoroutineNumber + 1
				continue
			}

			if CurrentGoroutineNumber > 0 && ItemNumber >= MaxBuy {
				MissionChannel <- MaxBuy
				ItemNumber = ItemNumber - MaxBuy
				CurrentGoroutineNumber = CurrentGoroutineNumber + 1
				continue
			}

			if CurrentGoroutineNumber > 0 && ItemNumber < MaxBuy {
				break
			}
		}

		if CurrentGoroutineNumber == 0 && ItemNumber == 0 {
			stopLooper()
			break
		}
	}
}

func SnapUp(ctx echo.Context) error {

	var result Result

	// 目前已經販完畢了，相關的資源都已經關閉了
	if MaxBuy == 0 && ItemNumber == 0 {
		result.Code = StatusSellOut
		return ctx.JSON(http.StatusServiceUnavailable, result)
	}

	// 檢查目前時間，如果時間還沒有到，返回「目前還不能購買」狀態
	if !checkTimeAvailable() {
		result.Code = StatusNotAllowToSell
		return ctx.JSON(http.StatusServiceUnavailable, result)
	}

	// 還沒販售完畢，但是目前可以取得的任務列隊已經歸零了
	// 代表「短期販售完畢了」，返回販售完畢狀態
	if len(MissionChannel) == 0 {
		result.Code = StatusSellOut
		return ctx.JSON(http.StatusServiceUnavailable, result)
	}

	// 綁定參數
	var input Input
	if i, err := strconv.Atoi(ctx.FormValue("customer_id")); err != nil {
		result.Code = StatusParamError
		return ctx.JSON(http.StatusBadRequest, result)
	} else {
		input.CustomerID = int32(i)
	}

	if i, err := strconv.Atoi(ctx.FormValue("number")); err != nil {
		result.Code = StatusParamError
		return ctx.JSON(http.StatusBadRequest, result)
	} else {
		input.Number = int32(i)
	}

	// 從 MissionChannel 當中拿出，如果不為 0，啟動執行動作
	availableNumber := <-MissionChannel
	if availableNumber == 0 {
		result.Code = StatusSellOut
		return ctx.JSON(http.StatusServiceUnavailable, result)
	}

	// 如果此次要買的數量大於目前可購買數量，取消這次的購買
	// 將數字推入至 ResultChannel
	if input.Number > availableNumber {
		ResultChannel <- availableNumber
		result.Code = StatusOverBuy
		return ctx.JSON(http.StatusServiceUnavailable, result)
	}

	go HandleMission(input, availableNumber)

	result.Code = StatusSuccess
	return ctx.JSON(http.StatusOK, result)
}

func checkTimeAvailable() bool {

	currentTime := time.Now()

	if currentTime.After(AvailableTime) {
		return true
	}

	return false
}

func HandleMission(input Input, availableNumber int32) {

	// 將目前資料庫的產品庫存扣掉
	// 將這次的購買計入寫到資料庫當中
	var updateInput add.Input
	updateInput.CustomerID = input.CustomerID
	updateInput.Number = input.Number
	updateInput.ProductID = ProductID

	err := add.Exec(updateInput)
	if err != nil {
		fmt.Println("creating error %s \n", err)
		// 這邊可以做一些動作來通知買家「此次交易失敗」
		ResultChannel <- availableNumber
	}

	leftItemNumber := MaxBuy - input.Number
	ResultChannel <- leftItemNumber
}

// 停止相關資源，等待被回收
func stopLooper() {
	MaxBuy = 0
	ItemNumber = 0

	close(MissionChannel)
	close(ResultChannel)
}
