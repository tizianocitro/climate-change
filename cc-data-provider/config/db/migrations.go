package db

import (
	"github.com/blang/semver"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Migration struct {
	fromVersion   semver.Version
	toVersion     semver.Version
	migrationFunc func(sqlx.Ext, *DB) error
}

const MySQLCharset = "DEFAULT CHARACTER SET utf8mb4"

var migrations = []Migration{
	{
		fromVersion: semver.MustParse("0.0.0"),
		toVersion:   semver.MustParse("0.1.0"),
		migrationFunc: func(e sqlx.Ext, db *DB) error {
			if _, err := e.Exec(`
				CREATE TABLE IF NOT EXISTS CCDP_System (
					SKey VARCHAR(64) PRIMARY KEY,
					SValue VARCHAR(1024) NULL
				);
			`); err != nil {
				return errors.Wrapf(err, "failed creating table CCDP_System")
			}

			if _, err := e.Exec(`
				CREATE TABLE IF NOT EXISTS CCDP_Issue (
					ID TEXT PRIMARY KEY,
					Name TEXT NOT NULL,
					ObjectivesAndResearchArea TEXT
				);
			`); err != nil {
				return errors.Wrapf(err, "failed creating table CCDP_Issue")
			}

			if _, err := e.Exec(`
				CREATE TABLE IF NOT EXISTS CCDP_Outcome (
					IssueID TEXT NOT NULL REFERENCES CCDP_Issue(ID),
					ID TEXT NOT NULL,
					Outcome TEXT
				);
			`); err != nil {
				return errors.Wrapf(err, "failed creating table CCDP_Outcome")
			}

			if _, err := e.Exec(`
				CREATE TABLE IF NOT EXISTS CCDP_Role (
					IssueID TEXT NOT NULL REFERENCES CCDP_Issue(ID),
					ID TEXT NOT NULL,
					UserID TEXT,
					Roles TEXT
				);
			`); err != nil {
				return errors.Wrapf(err, "failed creating table CCDP_Role")
			}

			if _, err := e.Exec(`
				CREATE TABLE IF NOT EXISTS CCDP_Element (
					IssueID TEXT NOT NULL REFERENCES CCDP_Issue(ID),
					ID TEXT NOT NULL,
					Name TEXT NOT NULL,
					Description TEXT,
					OrganizationID TEXT NOT NULL,
					ParentID TEXT NOT NULL
				);
			`); err != nil {
				return errors.Wrapf(err, "failed creating table CCDP_Element")
			}

			if _, err := e.Exec(`
				CREATE TABLE IF NOT EXISTS CCDP_Attachment (
					IssueID TEXT NOT NULL REFERENCES CCDP_Issue(ID),
					ID TEXT NOT NULL,
					Attachment TEXT
				);
			`); err != nil {
				return errors.Wrapf(err, "failed creating table CCDP_Attachment")
			}

			return nil
		},
	},
	{
		fromVersion: semver.MustParse("0.1.0"),
		toVersion:   semver.MustParse("0.2.0"),
		migrationFunc: func(e sqlx.Ext, db *DB) error {
			// prior to v1.0.0, this migration was used to trigger the data migration from the kvstore
			return nil
		},
	},
}
