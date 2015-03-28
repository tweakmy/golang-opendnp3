package DNP3

import (

)

//    buf := make([]byte, 0, 4096) // big buffer
//    tmp := make([]byte, 256)     // using small tmo buffer for demonstrating
//    for {
//        n, err := conn.Read(tmp)
//        if err != nil {
//            if err != io.EOF {
//                fmt.Println("read error:", err)
//            }
//            break
//        }
//        //fmt.Println("got", n, "bytes.")
//        buf = append(buf, tmp[:n]...)
//
//    }
//    fmt.Println("total size:", len(buf))