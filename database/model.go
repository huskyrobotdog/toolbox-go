package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/huskyrobotdog/toolbox-go/id"
)

type Basic struct {
	ID       id.ID `gorm:"primaryKey;column:id;" json:"id"`
	CreateAt id.ID `gorm:"column:create_at;" json:"create_at"`
	UpdateAt id.ID `gorm:"column:update_at;" json:"update_at"`
}

func NewBasic() Basic {
	_id := id.ID(id.Instance.Generate().Int64())
	now := id.ID(time.Now().Unix())
	return Basic{
		ID:       _id,
		CreateAt: now,
		UpdateAt: now,
	}
}

type JsonValue json.RawMessage

func (j *JsonValue) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JsonValue(result)
	return err
}

func (j JsonValue) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

func ToJsonValue(v any) (JsonValue, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return JsonValue(json.RawMessage(bytes)), nil
}
