package git

import "testing"

func TestValidateRepoURL(t *testing.T) {
	ok := []string{
		"https://github.com/org/repo.git",
		"http://example.com/r.git",
		"git@github.com:org/repo.git",
		"ssh://git@host/repo.git",
	}
	for _, u := range ok {
		if err := validateRepoURL(u); err != nil {
			t.Fatalf("%q: %v", u, err)
		}
	}
	bad := []string{"-u./evil", "file:///etc/passwd", "/tmp/repo", ""}
	for _, u := range bad {
		if u == "" {
			continue
		}
		if err := validateRepoURL(u); err == nil {
			t.Fatalf("expected error for %q", u)
		}
	}
}

func TestValidateBranch(t *testing.T) {
	if err := validateBranch("main"); err != nil {
		t.Fatal(err)
	}
	if err := validateBranch("feature/x-1"); err != nil {
		t.Fatal(err)
	}
	if err := validateBranch("-main"); err == nil {
		t.Fatal("expected error for -main")
	}
	if err := validateBranch("main;rm"); err == nil {
		t.Fatal("expected error for metachar")
	}
}

func TestValidateSHA(t *testing.T) {
	if err := validateSHA("abc1234"); err != nil {
		t.Fatal(err)
	}
	if err := validateSHA("deadbeefdeadbeefdeadbeefdeadbeefdeadbeef"); err != nil {
		t.Fatal(err)
	}
	if err := validateSHA("-abc1234"); err == nil {
		t.Fatal("expected error")
	}
	if err := validateSHA("not-hex!"); err == nil {
		t.Fatal("expected error")
	}
	if err := validateSHA("abc"); err == nil {
		t.Fatal("expected error for short sha")
	}
}
