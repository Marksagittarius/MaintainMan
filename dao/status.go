package dao

import (
	"maintainman/logger"
	"maintainman/model"
)

func GetOrderByRepairer(id uint, current bool, offset uint) (orders []*model.Order, err error) {
	status := &model.Status{
		RepairerID: id,
		Current:    current,
	}
	statuses := []*model.Status{}

	if err = Filter("id", offset, 0).Preload("Order").Where(status).Find(&statuses).Error; err != nil {
		logger.Logger.Debugf("GetOrderByRepairerErr: %v\n", err)
		return
	}

	for _, status := range statuses {
		orders = append(orders, status.Order)
	}

	return
}

func NewStatus(status, repairer uint, operator uint) *model.Status {
	return &model.Status{
		Status:     status,
		RepairerID: repairer,
		Current:    true,
		BaseModel: model.BaseModel{
			CreatedBy: operator,
			UpdatedBy: operator,
		},
	}
}

// StatusWaiting 待维修
func StatusWaiting(operator uint) *model.Status {
	return NewStatus(1, 0, operator)
}

// StatusAccepted 已接单
func StatusAssigned(repairer, operator uint) *model.Status {
	return NewStatus(2, repairer, operator)
}

// StatusCompleted 已完成
func StatusCompleted(operator uint) *model.Status {
	return NewStatus(3, 0, operator)
}

// StatusReported 上报中
func StatusReported(operator uint) *model.Status {
	return NewStatus(4, 0, operator)
}

// StatusHold 挂单
func StatusHold(operator uint) *model.Status {
	return NewStatus(5, 0, operator)
}

// StatusCanceled 已取消
func StatusCanceled(operator uint) *model.Status {
	return NewStatus(6, 0, operator)
}

// StatusRejected 已拒绝
func StatusRejected(operator uint) *model.Status {
	return NewStatus(7, 0, operator)
}

// StatusAppraised 已评价
func StatusAppraised(operator uint) *model.Status {
	return NewStatus(8, 0, operator)
}
