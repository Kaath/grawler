package grawler

import (
	"reflect"
	"sync"
	"testing"
	httpmock "github.com/jarcoal/httpmock"
)

func TestSafeCounter_SafeCount(t *testing.T) {
	type fields struct {
		count int
		mutex sync.Mutex
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
				mutex: tt.fields.mutex,
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
		mutex sync.Mutex
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
				mutex: tt.fields.mutex,
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
		page string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "No Links",
			args: args{page: "lkajsf;akj lkjdslfkjwejfhoweiu fajsldn jashof awjesljaengkjaljgn we;isjf'awejgl kejrhlbaewk f;aeksjnv lkjergdlvaewks f/AWLDGJNL ERGBLA WEKDSFNV"},
			want: nil},
		{name: "One Links",
			args: args{page: "lkajsf;akj lkjdslfkjwejfhoweiu fajsldn jashof awjesljaengkjaljgn we;isjf'awejgl kejrhlbaewk f;aeksjnv lkjergdlvaewks f/AWhref=\"https://www.google.com\"LDGJNL ERGBLA WEKDSFNV"},
			want: []string{"https://www.google.com"}},
		{name: "Multiple Simple Links",
			args: args{page: "lkajsf;akj lkjdslfkjwehref=\"https://www.linux.org\"jfhoweiu fajsldn jashof awjesljaehref=\"https://www.monculq.surla.commode.com\"ngkjaljgn we;isjf'awejgl kejrhlbaewk f;aeksjnv lkjergdlvaewks f/AWhref=\"https://www.google.com\"LDGJNL ERGBLA WEKDSFNV"},
			want: []string{"https://www.linux.org", "https://www.monculq.surla.commode.com", "https://www.google.com"}},
		{name: "One link ends with space",
			args: args{page: "a;sdjf;adsalkdsjadgadkj;asdkg;aslkg;asdkgj;asdlkjflek lkejvaidjvhref=http://www.google.com paoeid jgpoaierjg valdsvnlakdszxvnsleidjgp aviszdvnalzsjdknv ljeagvnapzsdixkjv; ezkldjvaerikdbjv;sdolxvkj"},
			want: []string{"http://www.google.com"}},
		{name: "One link ends with \\",
			args: args{page: "lflfklhfdlkjfkjlhflheflhgerlhg;gweghttps://www.google.com\\ergwrg"},
			want: []string{"https://www.google.com"}},
		{name: "One link ends with >",
			args: args{page: "kaf;dfajfpadj f;sdjfpoaihparuhawrjgnv;awgh;dskvj;sdgja;ig;asodigj;asoidj;wesdhref=http://www.google.com>sldkfjsldkjf;akjf;wekjfa;sdhg;adsijg'awegjariodjbadkjva;sdjv"},
			want: []string{"http://www.google.com"}},
		{name: "Multiple links",
			args: args{page: "asfjlaksdjfkshref=\"https://news.kaathe.busuttil.onl\"fasdf;asdkfj;askdjhref=http://www.google.com\\\\fjlkjndhref=http://www.linux.org>kjfkweifhifjoeijfowfj;aefjhref=https://www.monculq.com "},
			want: []string{"https://news.kaathe.busuttil.onl", "http://www.google.com", "http://www.linux.org", "https://www.monculq.com"}},
		{name: "One link ends with ' ' with query",
			args: args{page: "kaf;dfajfpadj f;sdjfpoaihparuhawrjgnv;awgh;dskvj;sdgja;ig;asodigj;asoidj;wesdhref=http://www.google.com/search?q=truc+machin sldkfjsldkjf;akjf;wekjfa;sdhg;adsijg'awegjariodjbadkjva;sdjv"},
			want: []string{"http://www.google.com/search?q=truc+machin"}},
		{name: "One link ends with > with query",
			args: args{page: "kaf;dfajfpadj f;sdjfpoaihparuhawrjgnv;awgh;dskvj;sdgja;ig;asodigj;asoidj;wesdhref=http://www.google.com/search?q=truc+machin>sldkfjsldkjf;akjf;wekjfa;sdhg;adsijg'awegjariodjbadkjva;sdjv"},
			want: []string{"http://www.google.com/search?q=truc+machin"}},
		{name: "One link ends with \\",
			args: args{page: "kaf;dfajfpadj f;sdjfpoaihparuhawrjgnv;awgh;dskvj;sdgja;ig;asodigj;asoidj;wesdhref=http://www.google.com/search?q=truc+machin\\sldkfjsldkjf;akjf;wekjfa;sdhg;adsijg'awegjariodjbadkjva;sdjv"},
			want: []string{"http://www.google.com/search?q=truc+machin"}},
		{name: "One link ends with \"",
			args: args{page: "kaf;dfajfpadj f;sdjfpoaihparuhawrjgnv;awgh;dskvj;sdgja;ig;asodigj;asoidj;wesdhref=\"http://www.google.com/search?q=truc+machin\"sldkfjsldkjf;akjf;wekjfa;sdhg;adsijg'awegjariodjbadkjva;sdjv"},
			want: []string{"http://www.google.com/search?q=truc+machin"}},
		{name: "Multiple links with queries and fragments",
			args: args{page: "asfjlaksdjfkshref=\"https://news.kaathe.busuttil.onl/search?jsp=jstp&q=truc+machin\"fasdf;asdkfj;askdjhref=http://www.google.com#thing\\\\fjlkjndhref=http://www.linux.org/list?chose=3#thing>kjfkweifhifjoeijfowfj;aefjhref=https://www.monculq.com "},
			want: []string{"https://news.kaathe.busuttil.onl/search?jsp=jstp&q=truc+machin", "http://www.google.com#thing", "http://www.linux.org/list?chose=3#thing", "https://www.monculq.com"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := find_urls(tt.args.page); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("find_urls() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_crawl(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://www.google.com",
		httpmock.NewStringResponder(200, "<!Doctype html><html><a href=\"https://old.sites.google.com#from-new\">blep</a><a href=http://insecure-old.sites.google.com/query?ip=192.168.9.2>blop</a></html>"))

	httpmock.RegisterResponder("GET", "https://old.sites.google.com#from-new",
		httpmock.NewStringResponder(200, "<!Doctype><html>ok</html>"))
	httpmock.RegisterResponder("GET", "http://insecure-old.sites.google.com/query?ip=192.168.9.2",
		httpmock.NewStringResponder(200, "<!Doctype><html>ok</html>"))

	type args struct {
		url        string
		depth      int
		treatments []SaveFunc
	}
	tests := []struct {
		name string
		args args
	}{

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crawl(tt.args.url, tt.args.depth, tt.args.treatments)
		})
	}
}

func TestStartCrawl(t *testing.T) {
	type args struct {
		starts     []string
		treatments []SaveFunc
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Do some mocking
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StartCrawl(tt.args.starts, tt.args.treatments...)
		})
	}
}
