package input

import (
	"strings"
	"testing"
)

func TestParseCommandValidCaseForAcceptChangeByReviewerCommand(t *testing.T) {
	type TestCase struct {
		input           string
		expectedBotName string
	}

	list := []TestCase{
		{
			input:           "@bot r+",
			expectedBotName: "bot",
		},
		{
			input:           "@bot-bot r+",
			expectedBotName: "bot-bot",
		},

		{
			input:           "    @bot r+",
			expectedBotName: "bot",
		},

		{
			input:           "@bot        r+",
			expectedBotName: "bot",
		},

		{
			input: `@bot        r+



	`,
			expectedBotName: "bot",
		},
	}
	for _, testcase := range list {
		input := testcase.input

		ok, cmd := ParseCommand(input)
		if !ok {
			t.Errorf("input: `%v` should be ok", input)
			continue
		}

		v, ok := cmd.(*AcceptChangeByReviewerCommand)
		if !ok {
			t.Errorf("input: `%v` should be AcceptChangeByReviewerCommand", input)
			continue
		}

		expected := testcase.expectedBotName
		if actual := v.BotName(); actual != expected {
			t.Errorf("input: `%v` should be the expected bot (`%v`) name but `%v`", input, expected, actual)
			continue
		}
	}
}

func TestParseCommandValidCaseForAcceptChangeByOthersCommand(t *testing.T) {
	type TestCase struct {
		input    string
		sender   string
		expected []string
	}

	list := []TestCase{
		{
			input:    "@bot r=KoujiroFrau",
			expected: []string{"KoujiroFrau"},
		},
		{
			input:    "  @bot    r=KoujiroFrau  ",
			expected: []string{"KoujiroFrau"},
		},

		{
			input:    "@bot r=KoujiroFrau-a",
			expected: []string{"KoujiroFrau-a"},
		},
		{
			input:    "  @bot    r=KoujiroFrau-a ",
			expected: []string{"KoujiroFrau-a"},
		},
		{
			input:    "@bot r=me",
			sender:   "KoujiroFrau",
			expected: []string{"KoujiroFrau"},
		},

		{
			input:    "@bot r=KoujiroFrau,pipimi",
			expected: []string{"KoujiroFrau", "pipimi"},
		},
		{
			input:    "  @bot r=KoujiroFrau,pipimi   ",
			expected: []string{"KoujiroFrau", "pipimi"},
		},
		{
			input:    "  @bot r=KoujiroFrau,  pipimi   ",
			expected: []string{"KoujiroFrau", "pipimi"},
		},
		{
			input:    "  @bot r=KoujiroFrau ,  pipimi   ",
			expected: []string{"KoujiroFrau", "pipimi"},
		},
		{
			input:    "  @bot r= KoujiroFrau ,  pipimi   ",
			expected: []string{"KoujiroFrau", "pipimi"},
		},

		{
			input:    "@bot r=KoujiroFrau-a,pipimi-b",
			expected: []string{"KoujiroFrau-a", "pipimi-b"},
		},
		{
			input:    "  @bot r=KoujiroFrau-a,pipimi-b   ",
			expected: []string{"KoujiroFrau-a", "pipimi-b"},
		},
		{
			input:    "  @bot r=KoujiroFrau-a,   pipimi-b   ",
			expected: []string{"KoujiroFrau-a", "pipimi-b"},
		},
		{
			input:    "  @bot r=KoujiroFrau-a  ,   pipimi-b   ",
			expected: []string{"KoujiroFrau-a", "pipimi-b"},
		},
		{
			input:    "  @bot r= KoujiroFrau-a  ,   pipimi-b   ",
			expected: []string{"KoujiroFrau-a", "pipimi-b"},
		},
		{
			input:    "@bot r=me, pipimi",
			sender:   "KoujiroFrau",
			expected: []string{"KoujiroFrau", "pipimi"},
		},
	}
	for _, testcase := range list {
		input := testcase.input

		ok, cmd := ParseCommand(input)
		if !ok {
			t.Errorf("input: `%v` should be ok", input)
			continue
		}

		v, ok := cmd.(*AcceptChangeByOthersCommand)
		if !ok {
			t.Errorf("input: `%v` should be AcceptChangeByOthersCommand", input)
			continue
		}

		if len(v.Reviewer) != len(testcase.expected) {
			t.Errorf("input: `%v` should be the expected length (`%v`) but the acutual length is `%v`", input, len(testcase.expected), len(v.Reviewer))
			continue
		}

		for i, actual := range v.Reviewer {
			expected := testcase.expected[i]
			//TODO: Check sender name if `r=me` is used, this if expression is incomplete.
			if strings.Index(input, "me") != -1 && expected == testcase.sender {
				actual = testcase.sender
			}
			if actual != expected {
				t.Errorf("input: `%v` should be the expected (`%v`) but `%v`", input, expected, actual)
				continue
			}
		}
	}
}

func TestParseCommandValidCaseForAssignReviewerCommand(t *testing.T) {
	type TestCase struct {
		input    string
		expected []string
	}

	list := []TestCase{
		{
			input:    "r? @reviewer",
			expected: []string{"reviewer"},
		},
		{
			input:    "r? @reviewer-a",
			expected: []string{"reviewer-a"},
		},
		{
			input:    "  r? @reviewer  ",
			expected: []string{"reviewer"},
		},
		{
			input:    "   r? @reviewer-a   ",
			expected: []string{"reviewer-a"},
		},

		{
			input:    "@reviewer r?",
			expected: []string{"reviewer"},
		},
		{
			input:    "@reviewer-a r?",
			expected: []string{"reviewer-a"},
		},
		{
			input:    "   @reviewer  r? ",
			expected: []string{"reviewer"},
		},
		{
			input:    "    @reviewer-a   r?",
			expected: []string{"reviewer-a"},
		},

		{
			input:    "r? @reviewer @reviewer2",
			expected: []string{"reviewer", "reviewer2"},
		},
		{
			input:    "r? @reviewer-a @reviewer-b",
			expected: []string{"reviewer-a", "reviewer-b"},
		},
		{
			input:    "  r? @reviewer  @reviewer2",
			expected: []string{"reviewer", "reviewer2"},
		},
		{
			input:    "   r? @reviewer-a   @reviewer-b",
			expected: []string{"reviewer-a", "reviewer-b"},
		},
	}
	for _, testcase := range list {
		input := testcase.input

		ok, cmd := ParseCommand(input)
		if !ok {
			t.Errorf("input: `%v` should be ok", input)
			continue
		}

		v, ok := cmd.(*AssignReviewerCommand)
		if !ok {
			t.Errorf("input: `%v` should be AssignReviewerCommand", input)
			continue
		}

		if len(v.Reviewer) != len(testcase.expected) {
			t.Errorf("input: `%v` should be the expected length (`%v`) but the acutual length is `%v`", input, len(testcase.expected), len(v.Reviewer))
			continue
		}

		for i, expected := range testcase.expected {
			if actual := v.Reviewer[i]; actual != expected {
				t.Errorf("input: `%v` should be the expected (`%v`) but `%v`", input, expected, actual)
				continue
			}
		}
	}
}

func TestParseCommandValidCaseForCancelApprovedByReviewerCommand(t *testing.T) {
	type TestCase struct {
		input           string
		expectedBotName string
	}

	list := []TestCase{
		{
			input:           "@bot r-",
			expectedBotName: "bot",
		},
		{
			input:           "@bot-bot r-",
			expectedBotName: "bot-bot",
		},

		{
			input:           "    @bot r-",
			expectedBotName: "bot",
		},

		{
			input:           "@bot        r-",
			expectedBotName: "bot",
		},

		{
			input: `@bot        r-



	`,
			expectedBotName: "bot",
		},
	}
	for _, testcase := range list {
		input := testcase.input

		ok, cmd := ParseCommand(input)
		if !ok {
			t.Errorf("input: `%v` should be ok", input)
			continue
		}

		v, ok := cmd.(*CancelApprovedByReviewerCommand)
		if !ok {
			t.Errorf("input: `%v` should be CancelApprovedByReviewerCommand", input)
			continue
		}

		expected := testcase.expectedBotName
		if actual := v.BotName(); actual != expected {
			t.Errorf("input: `%v` should be the expected bot (`%v`) name but `%v`", input, expected, actual)
			continue
		}
	}
}

func TestParseCommandInvalidCase(t *testing.T) {
	input := []string{
		"Hello, I'm john.",
		"",
		"bot r+",
		"@bot",

		// r+
		"@bot r +",
		"@bot r r+",
		"@bot r+ r",
		" @ bot r+",
		" @ bot r +",
		`@bot
    r+`,

		// r-
		"@bot r -",
		"@bot r r-",
		"@bot r- r",
		" @ bot r-",
		" @ bot r -",
		`@bot
    r-`,

		// r=reviewer
		"@bot r=",
		"@bot r =a",
		"@bot r = a",
		"@bot r r=a",
		"@bot r=a r",
		" @ bot r=a",
		" @ bot r = a",
		" @ bot r =a",
		`@bot
    r=a`,

		// @reviewer r?
		"@bot r r?",
		"@bot r? r",
		"@bot r? @bot2",
		"@bot r ?",
		" @ bot r?",
		" @ bot r ? ",
		`@bot
    r?`,

		// r? @reviewer
		"r? r @bot",
		"r? @bot r",
		"r? @bot r @bot2",
		"r ? @bot",
		" r? @ bot",
		" r ? @ bot ",
		`r?
    @bot`,
	}
	for _, item := range input {
		if ok, _ := ParseCommand(item); ok {
			t.Errorf("%v should not be ok", item)
		}
	}
}
