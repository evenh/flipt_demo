package actions

import (
	"flipt_demo/models"
	"fmt"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/x/responder"
)

type ItemsResource struct {
	buffalo.Resource
}

func (v ItemsResource) scope(c buffalo.Context) *pop.Query {
	tx := c.Value("tx").(*pop.Connection)
	return tx.Q()
}

// List gets all Items. This function is mapped to the path
// GET /items
func (v ItemsResource) List(c buffalo.Context) error {
	items := &models.Items{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := v.scope(c).PaginateFromParams(c.Params())
	q = q.Order("created_at desc")

	// Retrieve all Items from the DB
	if err := q.All(items); err != nil {
		return err
	}

	// Make Items available inside the html template
	c.Set("items", items)

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.HTML("items/index.plush.html"))
}

// Show gets the data for one Item. This function is mapped to
// the path GET /items/{item_id}
func (v ItemsResource) Show(c buffalo.Context) error {
	// Allocate an empty Item
	item := &models.Item{}

	// To find the Item the parameter item_id is used.
	if err := v.scope(c).Find(item, c.Param("item_id")); err != nil {
		return c.Error(404, err)
	}

	// Make item available inside the html template
	c.Set("item", item)

	return c.Render(200, r.HTML("items/show.plush.html"))
}

// New renders the form for creating a new Item.
// This function is mapped to the path GET /items/new
func (v ItemsResource) New(c buffalo.Context) error {
	// Make item available inside the html template
	c.Set("item", &models.Item{})

	return c.Render(200, r.HTML("items/new.plush.html"))
}

// Create adds a Item to the DB. This function is mapped to the
// path POST /items
func (v ItemsResource) Create(c buffalo.Context) error {
	// Allocate an empty Item
	item := &models.Item{}

	// Bind item to the html form elements
	if err := c.Bind(item); err != nil {
		return err
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(item)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		// Make item available inside the html template
		c.Set("item", item)

		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("items/new.plush.html"))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Item was created successfully")

	// and redirect to the items index page
	return c.Redirect(302, "/items/%s", item.ID)
}

// Edit renders a edit form for a Item. This function is
// mapped to the path GET /items/{item_id}/edit
func (v ItemsResource) Edit(c buffalo.Context) error {
	// Allocate an empty Item
	item := &models.Item{}

	if err := v.scope(c).Find(item, c.Param("item_id")); err != nil {
		return c.Error(404, err)
	}

	// Make item available inside the html template
	c.Set("item", item)
	return c.Render(200, r.HTML("items/edit.plush.html"))
}

// Update changes a Item in the DB. This function is mapped to
// the path PUT /items/{item_id}
func (v ItemsResource) Update(c buffalo.Context) error {
	// Allocate an empty Item
	item := &models.Item{}

	if err := v.scope(c).Find(item, c.Param("item_id")); err != nil {
		return c.Error(404, err)
	}

	// Bind Item to the html form elements
	if err := c.Bind(item); err != nil {
		return err
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	verrs, err := tx.ValidateAndUpdate(item)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		// Make item available inside the html template
		c.Set("item", item)

		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the edit.html template that the user can
		// correct the input.
		res := responder.Wants("javascript", func(c buffalo.Context) error {
			return c.Render(422, r.JavaScript("items/edit.plush.js"))
		})
		res.Wants("html", func(c buffalo.Context) error {
			return c.Render(422, r.HTML("items/edit.plush.html"))
		})
		return res.Respond(c)
	}

	res := responder.Wants("javascript", func(c buffalo.Context) error {
		c.Set("item", item)
		return c.Render(200, r.JavaScript("items/edit.plush.js"))
	})
	res.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a success message
		c.Flash().Add("success", "Item was updated successfully")

		// and redirect to the items index page
		return c.Redirect(302, "/items/%s", item.ID)
	})
	return res.Respond(c)
}

// Destroy deletes a Item from the DB. This function is mapped
// to the path DELETE /items/{item_id}
func (v ItemsResource) Destroy(c buffalo.Context) error {
	// Allocate an empty Item
	item := &models.Item{}

	// To find the Item the parameter item_id is used.
	if err := v.scope(c).Find(item, c.Param("item_id")); err != nil {
		return c.Error(404, err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	if err := tx.Destroy(item); err != nil {
		return err
	}

	res := responder.Wants("javascript", func(c buffalo.Context) error {
		c.Set("item", item)
		return c.Render(200, r.JavaScript("items/destroy.plush.js"))
	})
	res.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a flash message
		c.Flash().Add("success", "Item was destroyed successfully")

		// Redirect to the items index page
		return c.Redirect(302, "/items")
	})

	return res.Respond(c)
}
