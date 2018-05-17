package main

import (
    "bufio"
    "compress/zlib"
    "fmt"
    "io"
    "os"
    "time"
)

var (
    inFile     string
    outFile    string
    decompress bool
    operation  string
)

var (
    reader *bufio.Reader
    writer *bufio.Writer
)

func main() {
    parseArgs()

    var count int64
    var err error

    start := time.Now()

    in, err := os.Open(inFile)
    if err != nil {
        panic(err)
    }
    defer in.Close()

    out, err := os.Create(outFile)
    if err != nil {
        panic(err)
    }
    defer out.Close()

    reader = bufio.NewReaderSize(in, 64*1024)
    writer = bufio.NewWriterSize(out, 1*1024*1024)
    defer writer.Flush()

    if decompress {
        count, err = runDecompress()
    } else {
        count, err = runCompress()
    }

    if err != nil {
        panic(err)
    }

    fmt.Printf(
        "%s %s -> %s complete(%d bytes in %.2f secs)\n",
        operation,
        os.Args[1],
        os.Args[2],
        count,
        time.Since(start).Seconds(),
    )
}

func runCompress() (int64, error) {
    zip := zlib.NewWriter(writer)

    defer writer.Flush()
    defer zip.Close()

    return io.Copy(zip, reader)
}

func runDecompress() (int64, error) {
    zip, err := zlib.NewReader(reader)
    if err != nil {
        panic(err)
    }
    defer zip.Close()

    return io.Copy(writer, zip)
}

func parseArgs() {
    if len(os.Args) < 3 {
        printUsage()
        os.Exit(1)
    }

    if len(os.Args) < 4 {
        inFile = os.Args[1]
        outFile = os.Args[2]
        operation = "Compress"
        return
    }

    if os.Args[1] == "-d" {
        decompress = true
        operation = "Decompress"
    }

    inFile = os.Args[2]
    outFile = os.Args[3]
}

func printUsage() {
    fmt.Println("Usage: zl [-d] <input file> <output file>")
    fmt.Println("\t -d : Decompress input file into output file (default mode: compress)")
}
