package model

import (
	"xorm.io/xorm"
)

type Device struct {
	Id    int64  `json:"id"`    // 设备编号
	UId   int64  `json:"uid"`   // 关联用户
	NATId string `json:"natId"` // NAT编号
	Name  string `json:"name"`  // 设备名称
}

type DeviceModel struct {
	DB *xorm.Engine
}

// 新增设备
func (model DeviceModel) AddDevice(uid int64, name, nat string) bool {
	device := &Device{
		UId:   uid,
		NATId: nat,
		Name:  name,
	}
	_, err := model.DB.Insert(device)
	if err == nil {
		return true
	} else {
		return false
	}
}

// 更新设备名称
func (model DeviceModel) UpdateName(id int64, name string) bool {
	device := &Device{
		Id: id,
	}
	has, err := model.DB.Get(device)
	if err != nil || !has {
		return false
	}
	_, err = model.DB.ID(device.Id).Update(device)
	return err == nil
}

// 删除设备
func (model DeviceModel) DelDevice(id int64) bool {
	device := &Device{
		Id: id,
	}
	_, err := model.DB.Delete(device)
	return err == nil
}

// 检查NAT编号
func (model DeviceModel) CheckNATId(nat string) bool {
	device := &Device{
		NATId: nat,
	}
	has, err := model.DB.Get(device)
	if err != nil {
		return true
	}
	return has
}

// 获取设备信息
func (model DeviceModel) GetDevice(id int64) (*Device, error) {
	device := &Device{
		Id: id,
	}
	has, err := model.DB.Get(device)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return device, nil
}

// 获取设备信息
func (model DeviceModel) NATGetDevice(id string) (*Device, error) {
	device := &Device{
		NATId: id,
	}
	has, err := model.DB.Get(device)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return device, nil
}

// 获取设备列表
func (model DeviceModel) GetDeviceList(uid int64, page, num int) ([]Device, error) {
	var list []Device
	var err error
	if uid != 1 {
		err = model.DB.Where("uid = ?", uid).Desc("id").Limit(num, (page-1)*num).Find(&list)
	} else {
		err = model.DB.Desc("id").Limit(num, (page-1)*num).Find(&list)
	}
	if err != nil {
		return nil, err
	}
	return list, nil
}
