package pager_test

// ReservedBytes() returns whatever was passed in when constructed
// PageSize() returns whatever was passed in when constructed
// Close() closes the file
// NewPager() opens the file
// GetPage() returns the correct page
// GetPage() returns a page of the right size
// Duplicate calls to GetPage() uses the cache