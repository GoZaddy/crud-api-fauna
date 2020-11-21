package controllers

import (
	"fmt"
	f "github.com/fauna/faunadb-go/v3/faunadb"
	"github.com/gin-gonic/gin"
	"github.com/gozaddy/crud-api-fauna/customerrors"
	"github.com/gozaddy/crud-api-fauna/database"
	"github.com/gozaddy/crud-api-fauna/models"
	"github.com/mitchellh/mapstructure"
	"net/http"
)

type controller struct{
	DB database.FaunaDB
}


func NewController(db database.FaunaDB) *controller{
	return &controller{db}
}

func (ctrl *controller) AddReadingItem(c *gin.Context) error {
	var req struct{
		Title string `json:"title" form:"title" binding:"required"`
		Link string `json:"link" form:"link" binding:"omitempty,url"`
		Type string `json:"type" form:"type" binding:"required"`
		Author string `json:"author" form:"author" binding:"required"`
	}

	err := c.ShouldBind(&req)
	if err != nil{
		return customerrors.NewAppError(http.StatusBadRequest, err.Error())
	}

	//generate item id
	id, err := ctrl.DB.NewID()
	if err != nil{
		return err
	}

	//create new reading item
	item := models.NewReadingItem(id, req.Title, req.Link, req.Type, req.Author)

	//save new reading item to FaunaDB
	val, err := ctrl.DB.AddDocument(database.ReadingItemCollection, item)
	if err != nil{
		return err
	}

	fmt.Println(val)
	c.JSON(200, gin.H{
		"message": "New reading item successfully created!",
	})

	return nil
}

func (ctrl *controller) GetAllReadingItems(c *gin.Context) error {
	itemType := c.Query("type")
	var match f.Expr

	if itemType == ""{
		match = f.Documents(f.Collection(database.ReadingItemCollection))
	} else{
		match = f.MatchTerm(f.Index(database.ReadingItemByTypeIndex), itemType)
	}
	var items []models.ReadingItem
	val, err := ctrl.DB.FaunaClient().Query(
		f.Map(
			f.Paginate(
				match,
			),
			f.Lambda("docRef", f.Select(f.Arr{"data"},f.Get(f.Var("docRef")))),
		),
	)

	fmt.Println()

	if err != nil{
		return err
	}

	err = val.At(f.ObjKey("data")).Get(&items)
	if err != nil{
		return err
	}

	c.JSON(200, gin.H{
		"message": "Reading items retrieved successfully",
		"data": items,
	})

	return nil
}

func (ctrl *controller) GetOneReadingItem(c *gin.Context) error {
	var item models.ReadingItem
	val, err := ctrl.DB.GetDocument(database.ReadingItemCollection, c.Param("id"))
	if err != nil{
		if ferr, ok := err.(f.FaunaError); !ok{
			return err
		} else{
			if ferr.Status() == 404 {
				return customerrors.NewAppError(ferr.Status(),"The reading item with the provided ID does not exist")
			} else if ferr.Status() == 400 {
				return customerrors.NewAppError(ferr.Status(), err.Error())
			}
			return err
		}
	}

	err = val.At(f.ObjKey("data")).Get(&item)
	if err != nil{
		return err
	}

	c.JSON(200, gin.H{
		"message": "Reading item retrieved!",
		"data": item,
	})
	return nil
}

func (ctrl *controller) UpdateOneReadingItem(c *gin.Context) error {
	var req struct{
		Title string `json:"title" form:"title" binding:"omitempty" fauna:"title" mapstructure:"title,omitempty"`
		Link string `json:"link" form:"link" binding:"omitempty,url" fauna:"link" mapstructure:"link,omitempty"`
		Type string `json:"type" form:"type" binding:"omitempty" fauna:"type" mapstructure:"type,omitempty"`
		Author string `json:"author" form:"author" binding:"omitempty" fauna:"author" mapstructure:"author,omitempty"`
	}

	err := c.ShouldBind(&req)
	if err != nil{
		return customerrors.NewAppError(http.StatusBadRequest, err.Error())
	}

	var update map[string]interface{}

	err = mapstructure.Decode(req, &update)
	if err != nil{
		return err
	}

	fmt.Println(update)

	err = ctrl.DB.UpdateDocument(database.ReadingItemCollection, c.Param("id"), update)
	if err != nil{
		if ferr, ok := err.(f.FaunaError); !ok{
			return err
		} else{
			if ferr.Status() == 404 {
				return customerrors.NewAppError(ferr.Status(),"The reading item with the provided ID does not exist")
			} else if ferr.Status() == 400 {
				return customerrors.NewAppError(ferr.Status(), err.Error())
			}
			return err
		}
	}

	c.JSON(200, gin.H{
		"message": "Reading item updated successfully!",
	})

	return nil
}

func (ctrl *controller) DeleteOneReadingItem(c *gin.Context) error {
	err := ctrl.DB.DeleteDocument(database.ReadingItemCollection, c.Param("id"))
	if err != nil{
		if ferr, ok := err.(f.FaunaError); !ok{
			return err
		} else{
			if ferr.Status() == 404 {
				return customerrors.NewAppError(ferr.Status(),"The reading item with the provided ID does not exist")
			} else if ferr.Status() == 400 {
				return customerrors.NewAppError(ferr.Status(), err.Error())
			}
			return err
		}
	}

	c.JSON(200, gin.H{
		"message": "Reading item deleted successfully",
	})

	return nil
}
