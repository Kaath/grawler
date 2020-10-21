package grawler

import (
	"reflect"
	"testing"

	httpmock "github.com/jarcoal/httpmock"
)

func TestSafeCounter_SafeCount(t *testing.T) {
	type fields struct {
		count int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{name: "Count 0", fields: fields{count: 0}, want: 0},
		{name: "Count 90", fields: fields{count: 90}, want: 90},
		{name: "Count -1", fields: fields{count: -1}, want: -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SafeCounter{
				count: tt.fields.count,
			}
			if got := s.SafeCount(); got != tt.want {
				t.Errorf("SafeCounter.SafeCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSafeCounter_SafeInc(t *testing.T) {
	type fields struct {
		count int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{name: "Increase 0", fields: fields{count: 0}, want: 1},
		{name: "Increase 90", fields: fields{count: 90}, want: 91},
		{name: "Increase -1", fields: fields{count: -1}, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SafeCounter{
				count: tt.fields.count,
			}
			s.SafeInc()
			if got := s.count; got != tt.want {
				t.Errorf("s.Count = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_find_urls(t *testing.T) {
	type args struct {
		page *Page
		pol  Policy
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "No Links",
			args: args{
				page: &Page{
					url: "http://placeholder.com",
					body: []byte("lkajsf;akj lkjdslfkjwejfhoweiu fajsldn jashof awjesljaengkjaljgn we;isjf'awejgl kejrhlbaewk f;aeksjnv lkjergdlvaewks f/AWLDGJNL ERGBLA WEKDSFNV"),
				},
				pol: ACCEPT_ALL,
			},
			want: nil,
		},
		{name: "One link",
			args: args{
				page: &Page{
					url: "http://placeholder.com",
					body: []byte("lkajsf;akj lkjdslfkjwejfhoweiu fajsldn jashof awjesljaengkjaljgn we;isjf'awejgl kejrhlbaewk f;aeksjnv lkjergdlvaewks f/AWhref=\"https://www.google.com\"LDGJNL ERGBLA WEKDSFNV"),
				},
				pol: ACCEPT_ALL,
			},
			want: []string{"https://www.google.com"},
		},
		{name: "Multiple Simple links",
			args: args{
				page: &Page{
					url: "http://placeholder.com",
					body: []byte("lkajsf;akj lkjdslfkjwehref=\"https://www.linux.org\"jfhoweiu fajsldn jashof awjesljaehref=\"https://www.monculq.surla.commode.com\"ngkjaljgn we;isjf'awejgl kejrhlbaewk f;aeksjnv lkjergdlvaewks f/AWhref=\"https://www.google.com\"LDGJNL ERGBLA WEKDSFNV"),
				},
				pol: ACCEPT_ALL,
			},
			want: []string{"https://www.linux.org", "https://www.monculq.surla.commode.com", "https://www.google.com"},
		},
		{name: "One link ends with space",
			args: args{
				page: &Page{
					url: "http://placeholder",
					body: []byte("a;sdjf;adsalkdsjadgadkj;asdkg;aslkg;asdkgj;asdlkjflek lkejvaidjvhref=http://www.google.com paoeid jgpoaierjg valdsvnlakdszxvnsleidjgp aviszdvnalzsjdknv ljeagvnapzsdixkjv; ezkldjvaerikdbjv;sdolxvkj"),
				},
				pol: ACCEPT_ALL,
			},
			want: []string{"http://www.google.com"},
		},
		{name: "One link ends with space",
			args: args{
				page: &Page{
					url: "http://placeholder",
					body: []byte("a;sdjf;adsalkdsjadgadkj;asdkg;aslkg;asdkgj;asdlkjflek lkejvaidjvhref=http://www.google.com paoeid jgpoaierjg valdsvnlakdszxvnsleidjgp aviszdvnalzsjdknv ljeagvnapzsdixkjv; ezkldjvaerikdbjv;sdolxvkj"),
				},
				pol: ACCEPT_ALL,
			},
			want: []string{"http://www.google.com"},
		},
		{name: "One link ends with \\",
			args: args{
				page: &Page{
					url: "Placeholder",
					body: []byte("lflfklhfdlkjfkjlhflheflhgerlhg;gweghttps://www.google.com\\ergwrg"),
				},
				pol: ACCEPT_ALL,
			},
			want: []string{"https://www.google.com"},
		},
		{name: "One link ends with >",
			args: args{
				page: &Page{
					url: "Placeholder",
					body: []byte("kaf;dfajfpadj f;sdjfpoaihparuhawrjgnv;awgh;dskvj;sdgja;ig;asodigj;asoidj;wesdhref=http://www.google.com>sldkfjsldkjf;akjf;wekjfa;sdhg;adsijg'awegjariodjbadkjva;sdjv"),
				},
				pol: ACCEPT_ALL,
			},
			want: []string{"http://www.google.com"},
		},
		{name: "Multiple links",
			args: args{
				page: &Page{
					url: "Placeholder",
					body: []byte("asfjlaksdjfkshref=\"https://news.kaathe.busuttil.onl\"fasdf;asdkfj;askdjhref=http://www.google.com\\\\fjlkjndhref=http://www.linux.org>kjfkweifhifjoeijfowfj;aefjhref=https://www.monculq.com "),
				},
				pol: ACCEPT_ALL,
			},
			want: []string{"https://news.kaathe.busuttil.onl", "http://www.google.com", "http://www.linux.org", "https://www.monculq.com"},
		},
		{name: "One link ends with ' ' with query",
			args: args{
				page: &Page{
					url: "Placeholder",
					body: []byte("kaf;dfajfpadj f;sdjfpoaihparuhawrjgnv;awgh;dskvj;sdgja;ig;asodigj;asoidj;wesdhref=http://www.google.com/search?q=truc+machin sldkfjsldkjf;akjf;wekjfa;sdhg;adsijg'awegjariodjbadkjva;sdjv"),
				},
				pol: ACCEPT_ALL,
			},
			want: []string{"http://www.google.com/search?q=truc+machin"},
		},
		{name: "One link ends with > with query",
			args: args{
				page: &Page{
					url: "Placeholder",
					body: []byte("kaf;dfajfpadj f;sdjfpoaihparuhawrjgnv;awgh;dskvj;sdgja;ig;asodigj;asoidj;wesdhref=http://www.google.com/search?q=truc+machin>sldkfjsldkjf;akjf;wekjfa;sdhg;adsijg'awegjariodjbadkjva;sdjv"),
				},
				pol: ACCEPT_ALL,
			},
			want: []string{"http://www.google.com/search?q=truc+machin"},
		},
		{name: "One link ends with \\",
			args: args{
				page: &Page{
					url: "Placeholder",
					body: []byte("kaf;dfajfpadj f;sdjfpoaihparuhawrjgnv;awgh;dskvj;sdgja;ig;asodigj;asoidj;wesdhref=http://www.google.com/search?q=truc+machin\\sldkfjsldkjf;akjf;wekjfa;sdhg;adsijg'awegjariodjbadkjva;sdjv"),
				},
				pol: ACCEPT_ALL,
			},
			want: []string{"http://www.google.com/search?q=truc+machin"},
		},
		{name: "One link ends with \"",
			args: args{
				page: &Page{
					url: "Placeholder",
					body: []byte("kaf;dfajfpadj f;sdjfpoaihparuhawrjgnv;awgh;dskvj;sdgja;ig;asodigj;asoidj;wesdhref=\"http://www.google.com/search?q=truc+machin\"sldkfjsldkjf;akjf;wekjfa;sdhg;adsijg'awegjariodjbadkjva;sdjv"),
				},
				pol: ACCEPT_ALL,
			},
			want: []string{"http://www.google.com/search?q=truc+machin"},
		},
		{name: "Multiple links with queries and fragments",
			args: args{
				page: &Page{
					url: "Placeholder",
					body: []byte("asfjlaksdjfkshref=\"https://news.kaathe.busuttil.onl/search?jsp=jstp&q=truc+machin\"fasdf;asdkfj;askdjhref=http://www.google.com#thing\\\\fjlkjndhref=http://www.linux.org/list?chose=3#thing>kjfkweifhifjoeijfowfj;aefjhref=https://www.monculq.com "),
				},
				pol: ACCEPT_ALL,
			},
			want: []string{"https://news.kaathe.busuttil.onl/search?jsp=jstp&q=truc+machin", "http://www.google.com#thing", "http://www.linux.org/list?chose=3#thing", "https://www.monculq.com"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := find_urls(tt.args.page, tt.args.pol); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("find_urls() = %v, want %v", got, tt.want)
			}
		})
	}
}

//}

func TestNew(t *testing.T) {
	type args struct {
		starts     []string
		treatments []SaveFunc
	}
	tests := []struct {
		name string
		args args
		want *Crawler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.starts, tt.args.treatments...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCrawler_crawl(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://www.google.com",
		httpmock.NewStringResponder(200, "<!Doctype html><html><a href=\"https://old.sites.google.com#from-new\">blep</a><a href=http://insecure-old.sites.google.com/query?ip=192.168.9.2>blop</a></html>"))

	httpmock.RegisterResponder("GET", "https://old.sites.google.com#from-new",
		httpmock.NewStringResponder(200, "<!Doctype><html>ok</html>"))
	httpmock.RegisterResponder("GET", "http://insecure-old.sites.google.com/query?ip=192.168.9.2",
		httpmock.NewStringResponder(200, "<!Doctype><html>ok</html>"))

	type fields struct {
		Starters      []string
		Treatments    []SaveFunc
		LogLevel      int
		StartPolicies []Policy
		NodePolicies  []Policy
		LeafPolicies  []Policy
		DefaultPolicy Policy
	}
	type args struct {
		url        string
		depth      int
		treatments []SaveFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "No policy : greedy behaviour",
			fields: fields{
				Starters:      []string{},
				Treatments:    []SaveFunc{},
				LogLevel:      0,
				StartPolicies: nil,
				NodePolicies:  nil,
				LeafPolicies:  nil,
				DefaultPolicy: ACCEPT_ALL,
			},
			args: args{url: "https://www.google.com", depth: 2, treatments: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Crawler{
				Starters:      tt.fields.Starters,
				Treatments:    tt.fields.Treatments,
				counter:       SafeCounter{count: 0},
				LogLevel:      tt.fields.LogLevel,
				StartPolicies: tt.fields.StartPolicies,
				NodePolicies:  tt.fields.NodePolicies,
				LeafPolicies:  tt.fields.LeafPolicies,
				DefaultPolicy: tt.fields.DefaultPolicy,
			}
			c.crawl(tt.args.url, tt.args.depth, tt.args.treatments)
		})
	}
}

func TestCrawler_StartCrawl(t *testing.T) {
	type fields struct {
		Starters      []string
		Treatments    []SaveFunc
		LogLevel      int
		StartPolicies []Policy
		NodePolicies  []Policy
		LeafPolicies  []Policy
		DefaultPolicy Policy
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Crawler{
				Starters:      tt.fields.Starters,
				Treatments:    tt.fields.Treatments,
				counter:       SafeCounter{ count: 0 },
				LogLevel:      tt.fields.LogLevel,
				StartPolicies: tt.fields.StartPolicies,
				NodePolicies:  tt.fields.NodePolicies,
				LeafPolicies:  tt.fields.LeafPolicies,
				DefaultPolicy: tt.fields.DefaultPolicy,
			}
			c.StartCrawl()
		})
	}
}
