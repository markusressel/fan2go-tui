package ui

import (
	"fan2go-tui/internal/uistate"
)

const (
	// sessionNamespace is the namespace of the session of each fan2go daemon in the state, see session
	sessionNamespace = "session"
	// maxRememberedSessions is the number of fan2go daemons the session is remembered for
	maxRememberedSessions = 20
)

// session is what the user was looking at when the application was connected to a fan2go daemon the last time,
// restored when connecting to it again. The selections of the pages are remembered by the pages themselves, see
// selectionPersister.
type session struct {
	// Page is the page shown
	Page Page `json:"page,omitempty"`
}

// selectionPersister is a page that remembers its selection, see MainPage.PersistSession.
type selectionPersister interface {
	PersistSelection(store *uistate.Store, server string, restore bool)
}

// PersistSession remembers the shown page and the selections of the pages for the fan2go daemon at the address
// server in store, and restores the ones of the last run if restore is true. Must be called before Init.
func (mainPage *MainPage) PersistSession(store *uistate.Store, server string, restore bool) {
	mainPage.sessionStore = store
	mainPage.server = server
	store.SetLimit(sessionNamespace, maxRememberedSessions)

	for _, pagesPage := range mainPage.pagesMap.Iterator() {
		if persister, ok := pagesPage.(selectionPersister); ok {
			persister.PersistSelection(store, server, restore)
		}
	}

	if !restore {
		return
	}
	last, ok := uistate.Get[session](store, sessionNamespace, server)
	if !ok {
		return
	}
	if _, exists := mainPage.pagesMap.Get(last.Page); exists {
		mainPage.showPage(last.Page)
	}
}

// rememberSession saves the current session, see PersistSession.
func (mainPage *MainPage) rememberSession() {
	uistate.Set(mainPage.sessionStore, sessionNamespace, mainPage.server, session{Page: mainPage.page})
}
