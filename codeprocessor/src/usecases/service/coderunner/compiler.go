package coderunner

type Compiler int

const (
	UNKNOWN Compiler = iota
	GCC
	CLANG
	PYTHON
	NODEJS
)

func (c Compiler) CompileCommand() []string {
	switch c {
	case GCC:
		return []string{"gcc", "-o", runnerFilePath(buildFileName), runnerFilePath(c.SourseFileName())}
	case CLANG:
		return []string{"clang", "-o", runnerFilePath(buildFileName), runnerFilePath(c.SourseFileName())}
	case PYTHON, NODEJS:
		return nil
	default:
		return nil
	}
}

func (c Compiler) RunCommand() []string {
	switch c {
	case GCC, CLANG:
		return []string{runnerFilePath(buildFileName)}
	case PYTHON:
		return []string{"python3", runnerFilePath(c.SourseFileName())}
	case NODEJS:
		return []string{"node", runnerFilePath(c.SourseFileName())}
	default:
		return nil
	}
}

func (c Compiler) SourseFileName() string {
	switch c {
	case GCC, CLANG:
		return "source.c"
	case PYTHON:
		return "source.py"
	case NODEJS:
		return "source.js"
	default:
		return ""
	}
}

func (c Compiler) IsInterpreter() bool {
	return c.CompileCommand() == nil
}

func (c Compiler) String() string {
	switch c {
	case GCC:
		return "gcc"
	case CLANG:
		return "clang"
	case PYTHON:
		return "python3"
	case NODEJS:
		return "node"
	default:
		return ""
	}
}

func CompilerFromString(compileName string) Compiler {
	switch compileName {
	case GCC.String():
		return GCC
	case CLANG.String():
		return CLANG
	case PYTHON.String():
		return PYTHON
	case NODEJS.String():
		return NODEJS
	default:
		return UNKNOWN
	}
}
