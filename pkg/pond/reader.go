package pond

/*
   #include <stdio.h>   // 如果要调用C.free 一定要包含对应的头文件
   #include <stdlib.h>  // 此段注释与 import "C" 之间不能有空格
*/
import "C"
import (
	"github.com/nekohor/hororuen/pkg/paths"
	"github.com/nekohor/hororuen/pkg/converts"
	"log"
	"math"
	"sync"
	"syscall"
	"unsafe"
)

var (
	CURVE_DATA_MAX_NUM = 4000
)

type Reader struct {
	readFunc uintptr
	mutex *sync.Mutex
}

func NewReader() *Reader {
	reader := &Reader{}
	handle, err := syscall.LoadLibrary("ReadDCADLL.dll")
	if err != nil {
		log.Println(err)
		panic("err in Load DLL Library")
	}
	//defer syscall.FreeLibrary(handle)

	readFunc, err := syscall.GetProcAddress(handle, "ReadData")
	if err != nil {
		panic("err in GetProcAddress")
	}
	reader.readFunc = readFunc
	reader.mutex = new(sync.Mutex)

	return reader
}


func (reader *Reader) ReadData(dcaPath string, signalName string) (int, []float64) {

	reader.mutex.Lock()
	defer reader.mutex.Unlock()

	size := CURVE_DATA_MAX_NUM
	dataArray := make([]float64, size)
	if paths.IsExist(dcaPath) == true {

		//callArgDcaPath := uintptr(unsafe.Pointer(StringToINT8Ptr(dcaPath)))
		//callArgSignalName := uintptr(unsafe.Pointer(StringToINT8Ptr(signalName)))

		gbkDcaPath, err := converts.Utf8ToGbk([]byte(dcaPath))
		if err != nil {
			log.Println(err)
		}

		gbkSignalName, err := converts.Utf8ToGbk([]byte(signalName))
		if err != nil {
			log.Println(err)
		}

		CgoDcaPath := C.CString(string(gbkDcaPath))
		CgoSignalName := C.CString(string(gbkSignalName))
		defer C.free(unsafe.Pointer(CgoDcaPath))
		defer C.free(unsafe.Pointer(CgoSignalName))

		callArgDcaPath := uintptr(unsafe.Pointer(CgoDcaPath))
		callArgSignalName := uintptr(unsafe.Pointer(CgoSignalName))
		callArgDataArray := uintptr(unsafe.Pointer(&dataArray[0]))

		sizeUintptr, _, _ := syscall.Syscall(
			reader.readFunc, 3,
			callArgDcaPath,
			callArgSignalName,
			callArgDataArray)

		size = int(sizeUintptr)

	} else {
		log.Println("dcaPath does not exist: ", dcaPath)
		size = -2
	}

	log.Println("actual return size")
	log.Println(size)

	if -1 == size || -2 == size {
		size = 1
		log.Println("[Warning] wrong DCA path or signal name in DLL function")
	}

	// NaN转换为0
	buffArray := make([]float64, len(dataArray))
	for i := 0; i < len(dataArray); i++ {
		if math.IsNaN(dataArray[i]) {
			buffArray[i] = 0
		} else {
			buffArray[i] = dataArray[i]
		}
	}
	return size, buffArray
}