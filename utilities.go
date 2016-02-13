package main 


//generic error handling
func check(e error) {
    if e != nil {
        panic(e)
    }
}