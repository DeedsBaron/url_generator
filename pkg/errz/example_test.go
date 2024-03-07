package errz_test

import (
	"fmt"

	"gitlab.ozon.ru/express/core/lib/bx-common/errz"
)

var (
	catNotFound   = errz.NewCode("CatNotFound", errz.NotFound)
	notFoundDBErr = errz.NewCode("NotFoundDBErr", errz.Internal)
)

func Example() {
	_, err := whatIsCatNameServe()
	fmt.Println(err)

	_, err = getCatByIDServe()
	fmt.Println(err)

	// Output:
	// code = CatNotFound desc = failed to get cat by id: cat not found by id=1
	// code = NotFoundDBErr desc = failed to get cat by id: cat not found by id=1
}

// whatIsCatNameServe returns error: "code = CatNotFound desc = failed to get cat by id: cat not found by id=1"
//
// note: service will return NotFound as it was last wrapped error code
func whatIsCatNameServe() (string, error) {
	name, err := getCatByID(1)
	if err != nil {
		return "", errz.WrapC(err, catNotFound, "failed to get cat by id")
	}
	return name, nil
}

// getCatByIDServe returns error: "code = NotFoundDBErr desc = failed to get cat by id: cat not found by id=1"
func getCatByIDServe() (string, error) {
	name, err := getCatByID(1)
	if err != nil {
		return "", errz.Wrap(err, "failed to get cat by id")
	}
	return name, nil
}

func getCatByID(id int64) (string, error) {
	if id == 0 {
		return "myCat", nil
	}
	return "", errz.Errorf(notFoundDBErr, "cat not found by id=%v", id)
}
