package debug

import (
	"os"
	"runtime"
	"runtime/debug"
	"testing"
)

type Obj struct {
	x, y int
}

func objfin(x *Obj) {
	//println("finalized", x)
}

func TestWriteHeapDumpFinalizers(t *testing.T) {
	if runtime.GOOS == "js" {
		t.Skipf("WriteHeapDump is not available on %s.", runtime.GOOS)
	}
	f, err := os.CreateTemp("", "heapdumptest")
	if err != nil {
		t.Fatalf("TempFile failed: %v", err)
	}
	//defer os.Remove(f.Name())
	defer f.Close()

	// bug 9172: WriteHeapDump couldn't handle more than one finalizer
	println("allocating objects")
	x := &Obj{}
	runtime.SetFinalizer(x, objfin)
	y := &Obj{}
	runtime.SetFinalizer(y, objfin)

	// Trigger collection of x and y, queueing of their finalizers.
	println("starting gc")
	runtime.GC()

	// Make sure WriteHeapDump doesn't fail with multiple queued finalizers.
	println("starting dump")
	//os.Stdout
	debug.WriteHeapDump(f.Fd())
	println("done dump")
}
