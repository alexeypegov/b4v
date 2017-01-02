package controller

import (
  "github.com/alexeypegov/b4v/model"  
)

// Context holds handler context parameters
type Context struct {
	DB *model.DB
}
