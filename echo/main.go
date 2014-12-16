package main

// This is an example of Emboxen framework.
// Just returns the inputed source code as the application output.

import (
	"github.com/yvt/emboxen-framework"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	emboxen.Accept(acceptor)
}

func acceptor(srcReader emboxen.SourceCodeReader) {
	tmp, err := ioutil.TempFile("", "echo")
	if err != nil {
		srcReader.ReportInternalError(err)
		return
	}

	tmpPath := tmp.Name()
	defer tmp.Close()
	defer os.Remove(tmpPath)

	_, err = io.Copy(tmp, srcReader)
	if err != nil {
		srcReader.ReportInternalError(err)
		return
	}

	_, err = tmp.Seek(0, 0)
	if err != nil {
		srcReader.ReportInternalError(err)
		return
	}

	buildOutputWriter := srcReader.CloseAndOpenBuildOuputWriter()
	buildOutputWriter.Write([]byte("Now returning the inputed source code..."))

	appOutputWriter := buildOutputWriter.CloseAndOpenApplicationOutputWriter()

	_, err = io.Copy(appOutputWriter, tmp)
	if err != nil {
		srcReader.ReportInternalError(err)
		return
	}

	appOutputWriter.Close()
}
