package stub

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/pivotal-cf-experimental/mysql-operator/pkg/apis/binding/v1alpha1"

	_ "github.com/go-sql-driver/mysql"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

// Returns an int >= min, < max
func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

// Generate a random string of A-Z chars with len = l
func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	log.Printf("event = %+v\n", event)

	switch o := event.Object.(type) {
	case *v1alpha1.MysqlBinding:
		log.Printf("o.Spec.Username = %+v\n", o.Spec.Username)

		// create
		// 1. connect to mysql (need root creds)
		// 2. create a random strings for: username, password & databse
		// 3. create a user (random username string, random password string)
		// 4. create a database (random database string, owned by the user)
		// 5. set username, password, databasename on object
		// 6. set statusa on the object

		adminPassword := os.Getenv("MARIADB_ROOT_PASSWORD")
		adminHost := os.Getenv("SERVICE_NAME")

		if o.Spec.Username != "" {
			logrus.Println("This already exists", o.Spec.Username)
			return nil
		}

		connectionString := fmt.Sprintf("root:%s@tcp(%s:3306)/mysql", adminPassword, adminHost)
		logrus.Printf("connectionString = %+v\n", connectionString)

		db, err := sql.Open("mysql", connectionString)
		if err != nil {
			logrus.Errorf("Failed to connect to mysql : %v", err)
			return err
		}
		defer db.Close()

		var (
			randomUsername = randomString(10)
			randomPassword = randomString(10)
			randomDatabase = randomString(10)
		)

		db_ddl := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", randomDatabase)
		logrus.Printf("db_ddl = %+v\n", db_ddl)
		_, err = db.Exec(db_ddl)
		if err != nil {
			logrus.Errorf("Failed to run CREATE DATABASE : %v", err)
			return err
		}

		user_ddl := fmt.Sprintf("CREATE USER IF NOT EXISTS '%s'@'%%' IDENTIFIED BY '%s'", randomUsername, randomPassword)
		logrus.Printf("user_ddl = %+v\n", user_ddl)
		_, err = db.Exec(user_ddl)
		if err != nil {
			logrus.Errorf("Failed to run CREATE USER : %v", err)
			return err
		}

		grant_ddl := fmt.Sprintf("GRANT ALL ON %s.* TO '%s'@'%%'", randomDatabase, randomUsername)
		logrus.Printf("grant_ddl = %+v\n", grant_ddl)
		_, err = db.Exec(grant_ddl)
		if err != nil {
			logrus.Errorf("Failed to run GRANT ALL: %v", err)
			return err
		}

		o.Spec.Username = randomUsername
		o.Spec.Password = randomPassword
		o.Spec.Hostname = adminHost
		o.Spec.Database = randomDatabase

		sdk.Update(o)

		// ???
		// questions:
		// - do we want to, no matter what, ensure that users/dbs exist?
		// - do we want to support changing a password (or other things)?
		//   - how can we have fields the user is allowed to change (the password)
		//     and others the only the controller is allowed to change/set (the
		//     database name, the hostname, ...)?

		// delete
		// - delete database, user, ...
		// questions:
		// - what if we are not able to delete the db right away (e.g. mysql is
		//   down for maintenance) -- do we still get the delete event in the next
		//   reconcile loop?
		// - who deletes the object? does this codepath need to explicitely delete
		//   it or is it somehow magically deleted by the SDK
	}
	return nil
}

// newbusyBoxPod demonstrates how to create a busybox pod
func newbusyBoxPod(cr *v1alpha1.MysqlBinding) *corev1.Pod {
	labels := map[string]string{
		"app": "busy-box",
	}
	return &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "busy-box",
			Namespace: cr.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    "MysqlBinding",
				}),
			},
			Labels: labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}
