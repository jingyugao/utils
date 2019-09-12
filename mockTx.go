package utils

// import (
// 	"database/sql"
// 	"reflect"

// 	"github.com/jmoiron/sqlx"
// )

// type DbInterface interface {
// 	Rebind(query string) string
// 	Get(dest interface{}, query string, args ...interface{}) error
// 	Select(dest interface{}, query string, args ...interface{}) error
// 	Exec(query string, args ...interface{}) (sql.Result, error)
// 	NamedExec(query string, arg interface{}) (sql.Result, error)
// 	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
// }

// type MockDBs struct {
// 	rawDB    *sqlx.DB
// 	dbs      []DbInterface
// 	task     func(DbInterface, interface{}) error
// 	args     []interface{}
// 	mainCH   chan int
// 	chs      []chan int
// 	finished []bool
// 	errCH    chan error
// }

// func (mdbs *MockDBs) RunParallel() {

// 	for i := 0; i < 5; i++ {
// 		mdbs.rawDB.Exec("delete from redpacket_white_user_bankcards")
// 		mdbs.run()
// 	}
// 	close(mdbs.errCH)
// }

// func (mdbs *MockDBs) findOneOrder(order []int) int {
// 	size := len(order)
// 	s := make([]int, size+1)
// 	copy(s, order)
// 	for i := 0; i < len(mdbs.args); i++ {
// 		if mdbs.finished[i] {
// 			continue
// 		}
// 		s[size] = i
// 		if !mdbs.OrderDup(s) {
// 			return s[size]
// 		}
// 	}
// 	return -1
// }

// func (mdbs *MockDBs) OrderDup(order []int) bool {
// 	for _, o := range mdbs.orders {
// 		if reflect.DeepEqual(o, order) {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (mdbs *MockDBs) run() {
// 	mdbs.finished = make([]bool, len(mdbs.args))
// 	for i := range mdbs.args {
// 		go func(idx int) {
// 			mdbs.errCH <- mdbs.task(mdbs.dbs[idx], mdbs.args[idx])
// 			<-mdbs.chs[idx]
// 			mdbs.mainCH <- idx
// 		}(i)
// 	}
// 	order := []int{}
// 	mdbs.mainCH <- -1
// 	for {
// 		select {
// 		case i := <-mdbs.mainCH:
// 			if i != -1 {
// 				mdbs.finished[i] = true
// 			}
// 			idx := mdbs.findOneOrder(order)
// 			if idx == -1 {
// 				mdbs.orders = append(mdbs.orders, order)
// 				return
// 			}
// 			mdbs.chs[idx] <- 0
// 			order = append(order, idx)
// 		}
// 	}
// }

// func generateCases(n int) (cases [][]int) {
// 	s := []int{}
// 	for i := 0; i < n; i++ {
// 		s = append(s, i)
// 	}
// 	return
// }

// type GenMockDBs struct {
// 	task   func(DbInterface, interface{}) error
// 	args   []interface{}
// 	db     *sqlx.DB
// 	orders []int
// }

// func (gm *GenMockDBs) NewCase() {
// 	mdbs := NewMockDBs(gm.db, gm.task, gm.args)

// 	return
// }

// func NewMockDBs(db *sqlx.DB, task func(DbInterface, interface{}) error, args []interface{}) (mdb *MockDBs) {
// 	mainCH := make(chan int, 1)
// 	n := len(args)
// 	var chs []chan int
// 	var dbs []DbInterface
// 	for i := 0; i < n; i++ {
// 		idx := i
// 		chs = append(chs, make(chan int, 1))
// 		dbs = append(dbs, &mockDB{
// 			wait:   func() { <-chs[idx] },
// 			notify: func() { mainCH <- -1 },
// 			db:     db,
// 		})
// 	}
// 	mdb = &MockDBs{
// 		dbs:      dbs,
// 		chs:      chs,
// 		mainCH:   mainCH,
// 		errCH:    make(chan error, 1000),
// 		task:     task,
// 		args:     args,
// 		finished: make([]bool, n),
// 	}

// 	return
// }

// type mockDB struct {
// 	wait   func()
// 	notify func()
// 	db     DbInterface
// }

// func (m *mockDB) Rebind(query string) string {
// 	return m.db.Rebind(query)
// }

// func (m *mockDB) Get(dest interface{}, query string, args ...interface{}) error {
// 	m.wait()
// 	defer m.notify()
// 	return m.db.Get(dest, query, args...)
// }

// func (m *mockDB) Select(dest interface{}, query string, args ...interface{}) error {
// 	m.wait()
// 	defer m.notify()
// 	return m.db.Select(dest, query, args...)
// }

// func (m *mockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
// 	m.wait()
// 	defer m.notify()
// 	return m.db.Exec(query, args...)
// }

// func (m *mockDB) NamedExec(query string, arg interface{}) (sql.Result, error) {
// 	m.wait()
// 	defer m.notify()
// 	return m.db.NamedExec(query, arg)
// }

// func (m *mockDB) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
// 	m.wait()
// 	defer m.notify()
// 	return m.db.Queryx(query, args...)
// }
