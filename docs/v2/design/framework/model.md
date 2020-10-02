# Model

Model is an interface for data modelling. It's akin to client and server as a building block for services development. 
Where the server handles queries and the client makes queries, the model is used for storing and accessing data.

## Overview

The model builds on the store interface much like client/server build on transport/broker. The model is much like 
rails activerecord and provides a simple crud layer for accessing the underlying data store, allowing the developer 
to focus on their data types rather than raw storage.

## Design

Here's the proposed design for the model to be added as go-micro/model

### Interface 

```
type Model interface {
	// Initialise options
	Init(...Option) error
	// Retrieve options
	Options() Options
	// Register a type
	Register(v interface{}) error
	// Create a record
	Create(v interface, ...CreateOption) error
	// Read a record into v
	Read(v interface{}, ...ReadOption) error
	// Update a record
	Update(v interface{}, ...UpdateOption) error
	// Delete a record
	Delete(v interface{}, ...DeleteOption) error
	// Model implementation
	String() string
}
```

Additionally there is the potential to create an `Entity` value aligned with Message/Request in Client/Server that extracts the common 
values required for any entity to be stored.

```
type Entity interface {
	// Id of the record
	Id() string
	// Type of entity e.g User
	Type() string
	// Associated value
	Value() interface{}
}
```

This would expect the model to accept Entity in Create/Update/Delete rather than just an interface type. Similarly how client.Call 
accepts a client.Request or client.Publish accepts a client.Message.

### Storage

The model would accept go-micro/store.Store as an option for which it would use as storage. This allows the model to be agnostic 
of any storage type and keeps it pluggable much like the rest of the go-micro packages.

```
// passed as an option
m := model.NewModel(
	model.Store(store),
)
```

### Encoding

The model would also accept a codec which would allow simple marshaling/unmarshalling of data types much like the client/server. 
We quite simply pass codec.Marshaler which has the methods Marshal/Unmarshal to the model. The preference is to use the proto 
codec by default so that we have efficient serialisation of data especially with large data types.

```
m := model.NewModel(
	mode.Codec(codec),
)
```

## Usage

Usage would be as follows

Define your types, ideally in proto

```
message User {
  string name = 1;
  string email = 2;
}
```

Then register this against a model

```
// create a new model
m := mode.NewModel()

// register the User with the model
m.Register(new(User))

// register with index
m.Register(new(User), model.AddIndex("name"))
```

This should initialise the model to use the `users` table and map any necessary fields for indexing. Fields which need to be indexed 
are mapped as metadata in the store.Record. An alternative strategy would be to index everything or index string and int fields 
automatically. 


Once you have the model the access would be as follows

```
// create a record
m.Create(&User{Name: "john"})

// update a record
m.Update(&User{Name: "john", Email: "john@example.com"})
```

Considering we may need to store with a primary key/id the alternative form would be

```
entity := m.NewEntity(&User{Name: "john"})

// create the new entity
m.Create(entity)
```

##  Proto generation

With proto code generation this would evolve to the point where crud is generated for you much like for client/server

Assuming the following proto

```
service Users {
	...standard rpc
}

message User {
	string name = 1;
	string email = 2;
}
```

Then registering and using it with a service
```
// create a new service
service := micro.NewService()

// create a new user model
userModel := pb.NewUsersModel(service.Model())

// assuming a request comes in to update email the model handles updating the single field
userModel.UpdateEmail(req.User)
```

The proto generated code generation might look as follows

```
func NewUsersModel(m model.Model) (UsersModel, error) {
	if err := m.Register(new(User)); err != nil {
		return nil, err
	}
	return &usersModel{m}
}

type usersModel struct {
	model.Model
}

func (u *usersModel) UpdateEmail(user *User) error {
	return m.Update(user, model.UpdateField("email", user.Email))
}
```
