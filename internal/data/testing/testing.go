package testing

import (
	"context"
	// "fmt"
	"log"
	teEntity "testing/internal/entity/testing"

	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"

	"testing/pkg/errors"
	jaegerLog "testing/pkg/log"
)

type (
	// Data ...
	Data struct {
		db   *sqlx.DB
		stmt map[string]*sqlx.Stmt

		tracer opentracing.Tracer
		logger jaegerLog.Factory
	}

	// statement ...
	statement struct {
		key   string
		query string
	}
)

// Tambahkan query di dalam const
// getAllUser = "GetAllUser"
// qGetAllUser = "SELECT * FROM users"
const (
	getDataByName  = "GetDataByName"
	qGetDataByName = `SELECT * FROM testing_timdaen WHERE Nama= ?`

	getDataByCity  = "GetDataByCity"
	qGetDataByCity = `SELECT * FROM testing_timdaen WHERE Kota= ?`

	getDataByID  = "GetDataByID"
	qGetDataByID = `SELECT * FROM testing_timdaen WHERE UserID= ?`

	insertDataArray  = "InsertDataArray"
	qInsertDataArray = `INSERT INTO testing_timdaen
	(UserID, Nama, Usia, Kota, LastUpdate)
	VALUES (NULL, ?, ?, ?, CURRENT_TIMESTAMP)`

	editDataByID  = "EditDataByID"
	qEditDataByID = `UPDATE testing_timdaen
	SET 
        Nama = ?,
        Usia = ?,
        Kota = ?,
        LastUpdate = CURRENT_TIMESTAMP
    WHERE UserID = ?`
)

// Tambahkan query ke dalam key value order agar menjadi prepared statements
// readStmt = []statement{
// 	{getAllUser, qGetAllUser},
// }
var (
	readStmt = []statement{
		{getDataByName, qGetDataByName},
		{getDataByCity, qGetDataByCity},
		{getDataByID, qGetDataByID},
		{insertDataArray, qInsertDataArray},
		{editDataByID, qEditDataByID},
	}
)

// New ...
func New(db *sqlx.DB, tracer opentracing.Tracer, logger jaegerLog.Factory) Data {
	d := Data{
		db:     db,
		tracer: tracer,
		logger: logger,
	}

	d.initStmt()
	return d
}

func (d *Data) initStmt() {
	var (
		err   error
		stmts = make(map[string]*sqlx.Stmt)
	)

	for _, v := range readStmt {
		stmts[v.key], err = d.db.PreparexContext(context.Background(), v.query)
		if err != nil {
			log.Fatalf("[DB] Failed to initialize statement key %v, err : %v", v.key, err)
		}
	}

	d.stmt = stmts
}

// contoh implementasi ...
// func (d Data) GetShowname(ctx context.Context, movieID string) (string, error) {
// 	var (
// 		showname string
// 		err      error
// 	)

//// WAJIB ADA
// 	if span := opentracing.SpanFromContext(ctx); span != nil {
// 		span := d.tracer.StartSpan("SQL SELECT", opentracing.ChildOf(span.Context()))
// 		span.SetTag("mysql.server", "123.72.156.4")
// 		span.SetTag("mysql.database", "movie")
// 		span.SetTag("mysql.table", "showname")
// 		span.SetTag("mysql.query", "SELECT * FROM movie.showname WHERE movie_id="+movieID)
// 		defer span.Finish()
// 		ctx = opentracing.ContextWithSpan(ctx, span)
// 	}
//// WAJIB ADA

// 	// assumed data fetched from database
// 	showname = "Joni Bizarre Adventure"

//// OPTIONAL, DISARANKAN DIBUAT LOGGINGNYA
// 	d.logger.For(ctx).Info("SQL Query Success", zap.String("showname", showname))

//// WAJIB ADA, INI MERUPAKAN LOGGING KALAU TERJADI ERROR, BISA DIPASANG DI SERVICE DAN HANDLER, TIDAK HANYA DI DATA SAJA
// 	// if err != nil {
// 	// 	d.logger.For(ctx).Error("SQL Query Failed", zap.Error(err))
// 	// 	return showname, err
// 	// }
//// WAJIB ADA

// 	return showname, err
// }

func (d Data) GetDataByName(ctx context.Context, nama string) ([]teEntity.Testing, error) {
	var (
		rows  *sqlx.Rows
		user  teEntity.Testing
		users []teEntity.Testing
		err   error
	)

	rows, err = d.stmt[getDataByName].QueryxContext(ctx, nama)
	for rows.Next() {
		if err := rows.StructScan(&user); err != nil {
			return users, errors.Wrap(err, "[DATA][GetDataByName]")
		}
		users = append(users, user)
	}
	defer rows.Close()

	return users, err
}

func (d Data) GetDataByCity(ctx context.Context, kota string) ([]teEntity.Testing, error) {
	var (
		rows  *sqlx.Rows
		user  teEntity.Testing
		users []teEntity.Testing
		err   error
	)

	rows, err = d.stmt[getDataByCity].QueryxContext(ctx, kota)
	for rows.Next() {
		if err := rows.StructScan(&user); err != nil {
			return users, errors.Wrap(err, "[DATA][GetDataByCity]")
		}
		users = append(users, user)
	}
	defer rows.Close()

	return users, err
}

func (d Data) GetDataByID(ctx context.Context, userID int) (teEntity.Testing, error) {
	var (
		user teEntity.Testing
		err  error
	)
	
	if err := d.stmt[getDataByID].QueryRowxContext(ctx, userID).StructScan(&user); err != nil {
		return user, errors.Wrap(err, "[DATA][GetDataByID]")
	}

	return user, err
}

func (d Data) InsertDataArray(ctx context.Context, users teEntity.Testing) error {
	var err error

	if _, err = d.stmt[insertDataArray].ExecContext(ctx,
	users.Nama,
	users.Usia,
	users.Kota); err != nil {
		return errors.Wrap(err, "[DATA][InsertDataArray]")
	}
	
	return err
}

func (d Data) EditDataByID(ctx context.Context, users teEntity.Testing) error {
	var err error

	if _, err = d.stmt[editDataByID].ExecContext(ctx,
	users.Nama,
	users.Usia,
	users.Kota,
	users.UserID); err != nil {
		return errors.Wrap(err, "[DATA][EditDataByID]")
	}
	
	return err
}
