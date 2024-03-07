package tokengenerator_test

import (
	"fmt"
	"log"

	"github.com/AlexanderVasiliev23/yp-url-shortener/pkg/tokengenerator"
)

func Example() {
	const tokenLen = 8

	tokGen := tokengenerator.New(tokenLen)

	_token, err := tokGen.Generate()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(_token))

	// Output:
	// 8
}
