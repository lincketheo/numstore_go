package numstore

import "github.com/lincketheo/numstore/internal/nserror"

type Connection struct {
	DbName  string
	VarName string
}

func EmptyConnection() Connection {
  return Connection {
    DbName: "",
    VarName: "",
  }
}

func (c Connection) IsConnectedToDB() bool {
	return len(c.DbName) > 0
}

func (c Connection) IsConnectedToVar() bool {
	return len(c.VarName) > 0
}

func (c Connection) isValidState() bool {
	if c.IsConnectedToVar() && !c.IsConnectedToDB() {
		return false
	}
	return true
}

func (c Connection) connectionStateChange(changeState func() (Connection, error)) (Connection, error) {
	if !c.isValidState() {
		return c, nserror.Connection_InvalidState
	}

	newConn, err := changeState()
	if err != nil {
		return c, err
	}

	if !newConn.isValidState() {
		return c, nserror.Connection_InvalidState
	}

	return newConn, nil
}

func (c Connection) DisconnectDB() (Connection, error) {
	return c.connectionStateChange(func() (Connection, error) {
		if !c.IsConnectedToDB() {
			return c, nserror.Connection_DBNotConnected
		}
		return Connection{
			DbName:  "",
			VarName: "",
		}, nil
	})
}

func (c Connection) ConnectDB(dbName string) (Connection, error) {
	return c.connectionStateChange(func() (Connection, error) {
		if c.IsConnectedToDB() {
			return c, nserror.Connection_DBAlreadyConnected
		}
		if dbName == "" {
			return c, nserror.Connection_InvalidDBName
		}
		return Connection{
			DbName:  dbName,
			VarName: c.VarName,
		}, nil
	})
}

func (c Connection) ConnectVar(varName string) (Connection, error) {
	return c.connectionStateChange(func() (Connection, error) {
		if c.IsConnectedToVar() {
			return c, nserror.Connection_VarAlreadyConnected
		}
		if !c.IsConnectedToDB() {
			return c, nserror.Connection_DBNotConnected
		}
		if varName == "" {
			return c, nserror.Connection_InvalidVarName
		}
		return Connection{
			DbName:  c.DbName,
			VarName: varName,
		}, nil
	})
}
