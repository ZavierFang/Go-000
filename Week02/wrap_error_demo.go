package main

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"log"
)

type UserDao struct {
	uid  uint64
	name string
}

func main() {
	//调用路径 main -> aplService -> DomainService -> Dao
	aplService := &aplService{
		domainService: &domainService{
			userDao: &userDao{}}}

	_, err := aplService.findUser(context.Background(), 1)
	if err != nil {
		log.Printf("find user failed, err:%+v", err)
	}
}

//============AplService============
type aplService struct {
	domainService *domainService
}

func (p *aplService) findUser(ctx context.Context, uid uint64) (*UserDao, error) {
	//业务层也可以不处理往上抛error
	return p.domainService.findUserByUid(ctx, uid)
}

//============DomainService=========
type domainService struct {
	userDao *userDao
}

func (p *domainService) findUserByUid(ctx context.Context, uid uint64) (*UserDao, error) {
	userPo, err := p.userDao.findUserByUid(ctx, uid)
	if err != nil {
		//可以带上业务层处理的信息
		return nil, errors.WithMessage(err, "domain service: find user failed")
	}
	return userPo, nil
}

//============Dao===================
//将业务需要关心的错误映射成接口统一类型的错误返回
var errMap = map[error]error{
	sql.ErrNoRows: errors.New("no row in result"),
}

type userDao struct{}

func (p *userDao) findUserByUid(ctx context.Context, uid uint64) (*UserDao, error) {
	//something wrong before
	err := sql.ErrNoRows
	if repoErr, exist := errMap[err]; exist {
		err = repoErr
	}

	return nil, errors.Wrapf(err, "find user by uid(%d) failed", uid)
}
