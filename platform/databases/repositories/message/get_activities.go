package message

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"

)


func(r *repository)GetAllActive(ctx context.Context) ([]domain.Message,error){
 rows,err:=r.db.QueryContext(ctx,queryGetAllActive)
 if err!=nil{
	return nil,err
 }
 defer rows.Close()

 var messages []domain.Message
 for rows.Next(){
	var message domain.Message
	 err := rows.Scan(
			&message.ID,
			&message.Code,
			&message.Type,
			&message.Category,
			&message.Module,
			&message.Title,
			&message.Content,
			&message.Active,
			&message.CreatedAt,
			&message.UpdatedAt,
		)


	if err!=nil{
			return nil,err
	}
	messages=append(messages,message)
 }
 return messages,rows.Err()
}
