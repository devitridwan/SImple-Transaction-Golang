package service

import (
	"Ridwan/test_kanggo/src/model"
	"database/sql"
	"net/http"

	"github.com/labstack/echo"
)

func (this *ServiceKanggo) CreateProduct(context echo.Context, db *sql.DB) error {
	var req model.ReqCreateProduct
	var res = model.Response{Status: "failed"}
	if err := context.Bind(&req); err != nil {
		return context.JSON(http.StatusInternalServerError, res)
	} else {
		token := context.Request().Header.Get("Authorization")
		if val, ok := this.Token[token]; ok {
			var status bool
			err := db.QueryRow("select true from tbl_user where email = ? and status = ?", val, "admin").Scan(&status)
			if err == nil && status {
				_, err := db.Query("insert into tbl_product (name, price, qty) values(?, ?, ?)", req.Name,
					req.Price, req.Quantity)
				if err != nil {
					res.Error = err.Error()
					return context.JSON(http.StatusInternalServerError, res)
				} else {
					res.Status = "success"
				}
			}
		}

	}
	return context.JSON(http.StatusOK, res)
}

func (this *ServiceKanggo) EditProduct(context echo.Context, db *sql.DB) error {
	var req model.ReqUpdateProduct
	var res = model.Response{Status: "failed"}
	if err := context.Bind(&req); err != nil {
		return context.JSON(http.StatusInternalServerError, res)
	} else {
		token := context.Request().Header.Get("Authorization")
		if val, ok := this.Token[token]; ok {
			_, err := db.Query("update tbl_product set name=?, price=?, qty=? where id = ? and exists (select true from tbl_user where email = ? and status = ?)",
				req.Name, req.Price, req.Quantity, req.Id, val, "admin")
			if err != nil {
				res.Error = err.Error()
				return context.JSON(http.StatusInternalServerError, res)
			} else {
				res.Status = "success"
			}
		}

	}
	return context.JSON(http.StatusOK, res)
}

func (this *ServiceKanggo) ListProduct(context echo.Context, db *sql.DB) error {
	var res = model.ResProductList{Response: model.Response{Status: "failed"}}
	var data model.DataProductList
	token := context.Request().Header.Get("Authorization")
	if val, ok := this.Token[token]; ok {
		result, err := db.Query("select * from tbl_product where exists (select true from tbl_user where email = ?)", val)
		if err != nil {
			res.Error = err.Error()
			return context.JSON(http.StatusInternalServerError, res)
		} else {
			for result.Next() {
				if err := result.Scan(
					&data.Id,
					&data.Name,
					&data.Price,
					&data.Quantity); err != nil {
					res.Error = err.Error()
					return context.JSON(http.StatusInternalServerError, res)
				} else {
					res.Data = append(res.Data, data)
				}
			}
			res.Status = "success"
		}
	}

	// }
	return context.JSON(http.StatusOK, res)
}

func (this *ServiceKanggo) DeleteProduct(context echo.Context, db *sql.DB) error {
	var req model.ReqDeleteProduct
	var res = model.Response{Status: "failed"}
	if err := context.Bind(&req); err != nil {
		return context.JSON(http.StatusInternalServerError, res)
	} else {
		token := context.Request().Header.Get("Authorization")
		if val, ok := this.Token[token]; ok {
			_, err := db.Query("delete from tbl_product where id = ? and exists (select true from tbl_user where email = ? and status = ?)",
				req.Id, val, "admin")
			if err != nil {
				res.Error = err.Error()
				return context.JSON(http.StatusInternalServerError, res)
			} else {
				res.Status = "success"
			}
		}

	}
	return context.JSON(http.StatusOK, res)
}
