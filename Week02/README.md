## 课后作业

我们在数据库操作的时候，比如: dao层中当遇到一个 sql.ErrNoRows 的时候，是否应该Wrap这个error，抛给上层。为什么，应该怎么做请写出代码？

  需要Wrap后抛给上层。dao层属于基础层，基础层将错误抛给上层，由上层记录错误日志并转换为对外暴露的错误码或者做降级处理。
### 代码实现

**service层**
```go
type Service struct {
	dao *dao.Dao
}

func (svc *Service) QueryUserById(userID uint64) (*model.User, error) {
	user, err := svc.dao.QueryUserById(userID)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("user not found use mock data")
		user = model.User(Name: "xiaoming", Email: "xxx@xxx.com")
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
```

**dao层**
```go
type Dao struct {
	db *sqlx.DB
}

func NewDao(dsn string) (*Dao, error) {
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, errors.WithMessage(err, "connect error")
	}
	return &Dao{db: db}, nil
}

const queryUserBYId = "SELECT nick_name, mail FROM user WHERE id=? LIMIT 1"

func (d *Dao) QueryUserById(userID uint64) (*model.User, error) {
	var users []model.User
	err := d.db.Select(&users, queryUserBYId, userID)
	if err != nil {
		return nil, errors.Wrap(err, "query user by id")
	}
	return &users[0], nil
}

func (d *Dao) Close() {
	d.db.Close()
}
```