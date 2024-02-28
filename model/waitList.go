/*
* @Author: Oatmeal107
* @Date:   2023/6/19 17:33
 */

package model

import "gorm.io/gorm"

// WaitList 待审批记录, 上传的待审批记录是一批一批的,todo 先留着有需要再说
type WaitList struct {
	gorm.Model
}
