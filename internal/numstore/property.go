package numstore

type PropertyType interface {
	PrettyString() string
  Parse(json string) string
  Serialize() string
}

type PropertyTypeFactory interface {
  Default() PropertyType 
  Key() string
}

var propTypeRegistry = make(map[string]PropertyTypeFactory)

func Register(factory PropertyTypeFactory) {
  propTypeRegistry[factory.Key()] = factory
}


