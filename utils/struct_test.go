package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestStruct(t *testing.T) {
	type cat struct {
		Name string
	}
	ins := cat{Name: "mimi"}
	arr := GetStructKeys(ins, "json")
	if arr[0] != "Name" {
		t.Fail()
	}
	type tag struct {
		Tag int `json:"tag" id:"100"`
	}
	arr = GetStructKeys(tag{Tag: 1}, "json")
	if arr[0] != "tag" {
		t.Fail()
	}
}

func TestGetXormStructTag(t *testing.T) {
	type captcha struct {
		Id        uint64     `json:"id,string" xorm:"'id' pk autoincr" db:"'id'"`                       // 验证码id
		Code      string     `json:"code" xorm:"'code'" db:"'code'"`                                    // 验证码
		Target    string     `json:"target" xorm:"'target'" db:"'target'"`                              // 验证码发送对象
		Type      int16      `json:"type" xorm:"'type'" db:"'type'"`                                    // 验证码类型
		Scene     int16      `json:"scene" xorm:"'scene'" db:"'scene'"`                                 // 验证码场景
		CreatedAt *time.Time `json:"createdAt,omitempty" xorm:"'created_at' created" db:"'created_at'"` // 创建时间
		CreatedBy int64      `json:"createdBy,string" xorm:"'created_by'" db:"'created_by'"`            // 创建人
		ExpiredAt time.Time  `json:"expiredAt" xorm:"'expired_at'" db:"'expired_at'"`                   // 删除时间
	}
	{
		arr := GetXormStructKeys(captcha{})
		if len(arr) != 8 {
			t.Fail()
		}
		if arr[0] != "id pk autoincr" {
			t.Fail()
		}
	}
	{
		arr := GetDbStructKeys(captcha{})
		if len(arr) != 8 {
			t.Fail()
		}
		if arr[0] != "id" {
			t.Fail()
		}
	}
	{
		type LibraryFolder struct {
			FolderId      int64      `json:"folderId,string" xorm:"'folder_id' pk" db:"'folder_id'"`            // 主键ID
			FolderName    string     `json:"folderName" xorm:"'folder_name'" db:"'folder_name'"`                // 文件名称
			FolderData    string     `json:"folderData" xorm:"'folder_data'" db:"'folder_data'"`                //
			FolderLevel   int8       `json:"folderLevel" xorm:"'folder_level'" db:"'folder_level'"`             // 文件夹、回收站
			FolderFlags   int64      `json:"folderFlags,string" xorm:"'folder_flags'" db:"'folder_flags'"`      // falgs位
			FolderType    int8       `json:"folderType" xorm:"'folder_type'" db:"'folder_type'"`                // 文件夹类型
			FolderPath    string     `json:"folderPath" xorm:"'folder_path'" db:"'folder_path'"`                // 路径ID以/分割
			TmpId         int64      `json:"tmpId,string" xorm:"'tmp_id'" db:"'tmp_id'"`                        // 临时操作ID
			LibraryId     int64      `json:"libraryId,string" xorm:"'library_id'" db:"'library_id'"`            // 应用id
			AppId         int64      `json:"appId,string" xorm:"'app_id'" db:"'app_id'"`                        // 应用id
			CreatedBy     int64      `json:"createdBy,string" xorm:"'created_by'" db:"'created_by'"`            // 创建人id
			CreatedAt     *time.Time `json:"createdAt,omitempty" xorm:"'created_at' created" db:"'created_at'"` // 创建时间
			UpdatedAt     *time.Time `json:"updatedAt,omitempty" xorm:"'updated_at' updated" db:"'updated_at'"` // 更新时间
			DeletedAt     *time.Time `json:"deletedAt,omitempty" xorm:"'deleted_at'" db:"'deleted_at'"`         // 删除时间
			DeletedStatus int8       `json:"deletedStatus" xorm:"'deleted_status'" db:"'deleted_status'"`       // 删除状态
		}
		fmt.Println(GetDbStructKeys(LibraryFolder{}))
	}
}
