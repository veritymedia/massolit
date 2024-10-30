/**
* This file was @generated using pocketbase-typegen
*/

import type PocketBase from 'pocketbase'
import type { RecordService } from 'pocketbase'

export enum Collections {
	BookInstances = "book_instances",
	Books = "books",
	Keys = "keys",
	Rentals = "rentals",
	Users = "users",
}

// Alias types for improved usability
export type IsoDateString = string
export type RecordIdString = string
export type HTMLString = string

// System fields
export type BaseSystemFields<T = never> = {
	id: RecordIdString
	created: IsoDateString
	updated: IsoDateString
	collectionId: string
	collectionName: Collections
	expand?: T
}

export type AuthSystemFields<T = never> = {
	email: string
	emailVisibility: boolean
	username: string
	verified: boolean
} & BaseSystemFields<T>

// Record types for each collection

export type BookInstancesRecord = {
	book?: RecordIdString
	book_code: string
}

export type BooksRecord = {
	cover_url?: string
	isbn: number
	title: string
}

export type KeysRecord = {
	app?: string
	key?: string
}

export type RentalsRecord = {
	book_instance?: RecordIdString
	rented_to?: string
}

export type UsersRecord = {
	avatar?: string
	name?: string
}

// Response types include system fields and match responses from the PocketBase API
export type BookInstancesResponse<Texpand = unknown> = Required<BookInstancesRecord> & BaseSystemFields<Texpand>
export type BooksResponse<Texpand = unknown> = Required<BooksRecord> & BaseSystemFields<Texpand>
export type KeysResponse<Texpand = unknown> = Required<KeysRecord> & BaseSystemFields<Texpand>
export type RentalsResponse<Texpand = unknown> = Required<RentalsRecord> & BaseSystemFields<Texpand>
export type UsersResponse<Texpand = unknown> = Required<UsersRecord> & AuthSystemFields<Texpand>

// Types containing all Records and Responses, useful for creating typing helper functions

export type CollectionRecords = {
	book_instances: BookInstancesRecord
	books: BooksRecord
	keys: KeysRecord
	rentals: RentalsRecord
	users: UsersRecord
}

export type CollectionResponses = {
	book_instances: BookInstancesResponse
	books: BooksResponse
	keys: KeysResponse
	rentals: RentalsResponse
	users: UsersResponse
}

// Type for usage with type asserted PocketBase instance
// https://github.com/pocketbase/js-sdk#specify-typescript-definitions

export type TypedPocketBase = PocketBase & {
	collection(idOrName: 'book_instances'): RecordService<BookInstancesResponse>
	collection(idOrName: 'books'): RecordService<BooksResponse>
	collection(idOrName: 'keys'): RecordService<KeysResponse>
	collection(idOrName: 'rentals'): RecordService<RentalsResponse>
	collection(idOrName: 'users'): RecordService<UsersResponse>
}
