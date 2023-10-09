package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ws_comparator/domain/dto"
	"ws_comparator/infrastructure/persistence"
	"ws_comparator/infrastructure/repository"

	"github.com/go-resty/resty/v2"
	"github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
)

type ComparatorServiceImpl struct {
	comparatorRepository repository.ComparatorRepository
}

func InitComparatorServiceImpl() *ComparatorServiceImpl {
	dbHelper, err := persistence.InitDbHelper()
	if err != nil {
		log.Fatal(err.Error())
	}
	return &ComparatorServiceImpl{
		comparatorRepository: dbHelper.ComparatorRepository,
	}
}

func (r *ComparatorServiceImpl) Comparator(comparator dto.ComparatorIn) dto.Response {

	var response dto.Response
	/*configs, err := r.comparatorRepository.GetConfiguration()

	if err != nil {
		response.Status = http.StatusBadRequest
	}*/

	fmt.Println(comparator)

	payload := map[string]interface{}{}
	result := map[string]interface{}{}

	//resultdiff := map[string]interface{}{}

	for k, b := range comparator.Body {
		payload[k] = b
	}

	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post("https://coverage-microservice-qndxoltwga-uc.a.run.app/coverage.CoverageService/GetRoute")

	if err != nil {
		fmt.Println("Error realizando la solicitud:", err)
		//return err
	}

	fmt.Println("Estado de la respuesta:", resp.Status())

	err = json.Unmarshal(resp.Body(), &result)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result)

	//jsonDataResult, err := json.Marshal(result)

	//jsonDataIn, err := json.Marshal(comparator.ResponseC)

	/*	diffe := difflib.UnifiedDiff{
			A:        difflib.SplitLines(string(jsonDataResult)),
			B:        difflib.SplitLines(string(jsonDataIn)),
			FromFile: "Original",
			ToFile:   "Current",
			Context:  3,
		}
		text, _ := difflib.GetUnifiedDiffString(diffe)
		fmt.Printf(text)

		err = json.Unmarshal([]byte(text), &resultdiff)*/

	differ := gojsondiff.New()
	diff := differ.CompareObjects(result, comparator.ResponseC)
	if err != nil {
		fmt.Println("Error al comparar JSON:", err)
		//return
	}

	f := formatter.NewAsciiFormatter(result, formatter.AsciiFormatterDefaultConfig)
	deltaJson, _ := f.Format(diff)

	fmt.Println(deltaJson)

	/*for _, v := range diff.Deltas() {
		for _, r := range result {
			if r == string(v.Similarity()) {

			}
		}

		//resultdiff[v] =
	}*/

	/*	differ := gojsondiff.New()

		diff, err = differ.Compare(result, comparator.ResponseC)
		if err != nil {
			fmt.Println("Error al comparar JSON:", err)
			return
		}

		// Formatea las diferencias
		formatter := formatter.NewAsciiFormatter(originalMap, formatter.AsciiFormatterConfig{
			Coloring: true,
		})

		diffString, err := formatter.Format(diff)
		if err != nil {
			fmt.Println("Error al formatear diferencias:", err)
			return
		}

		// Imprime las diferencias
		fmt.Println("Diferencias entre JSON original y modificado:")
		fmt.Println(diffString)*/

	/*if diff = cmp.Diff(result, comparator.ResponseC); diff != "" {
		fmt.Printf("Las estructuras son diferentes:\n%s\n", diff)
	} else {
		fmt.Println("Las estructuras son iguales.")
	}*/

	//foodsDto := dto.FoodDto{}.TransformListEntityToDto(foods)
	/*if ok, resp := utils.ValidateQueryError(err,
		"NO_ROLES_FOUND", "ERROR_GETTING_ROLES", headers.Language); !ok {
		return userRoles, resp
	}

	result, err := r.relationUserRepository.FindUsersByRolAndTenantAndServiceDesk(
		rol.IdRol,
		headers.TenantId,
		idServiceDesk,
	)

	if ok, resp := utils.ValidateQueryError(err,
		"NO_USERS_FOUND", "ERROR_GETTING_USERS", headers.Language); !ok {
		return userRoles, resp
	}

	if len(result) <= 0 {
		response.Status = http.StatusNotFound
		response.Message = utils.Language(headers.Language, "NO_USERS_FOUND")
		return userRoles, response
	}

	userRoles.NameRol = rol.RolName

	for _, user := range result {
		userRoles.Users = append(userRoles.Users, dto.UserDto{
			IdUser:       user.IdUser,
			FullNameUser: user.UserName,
			Email:        user.Email,
		})
	}*/

	response.Status = http.StatusOK
	response.Description = diff.Deltas()
	return response
}
