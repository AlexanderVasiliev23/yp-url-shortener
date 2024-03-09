/*
Анализатор кода

 1. Добавлены все стандартные проверки golang.org/x/tools/go/analysis/passes

 2. Добавлены все анализаторы класса SA пакета staticcheck.io

 3. Добавлены по одному анализатору классов ST, S и QF пакета staticcheck.io

 4. Добавлены два анализатора из публичных github репозиториев:

    - rowserrcheck проверяет, что при запросах в бд ошибки sql.Rows.Err обрабатываются корректно

    - bodyclose проверяет, что не забыли закрыть resp.Body после получения ответа при http запросе

 5. Добавлен самописный анализатор osexitchecker, который проверяет, что в пакете main в функции main напрямую не вызывается функция os.Exit()
*/
package main

import (
	"github.com/AlexanderVasiliev23/yp-url-shortener/pkg/osexitchecker"
	"github.com/jingyugao/rowserrcheck/passes/rowserr"
	"github.com/timakin/bodyclose/passes/bodyclose"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/slog"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

func main() {
	var allChecks []*analysis.Analyzer

	allChecks = append(allChecks, checksFromAnalysis()...)
	allChecks = append(allChecks, staticcheckAnalyzers()...)
	allChecks = append(allChecks, otherPublicAnalyzers()...)
	allChecks = append(allChecks, osExitChecker())

	multichecker.Main(allChecks...)
}

func checksFromAnalysis() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		atomicalign.Analyzer,
		bools.Analyzer,
		buildssa.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		ctrlflow.Analyzer,
		deepequalerrors.Analyzer,
		defers.Analyzer,
		directive.Analyzer,
		errorsas.Analyzer,
		fieldalignment.Analyzer,
		findcall.Analyzer,
		framepointer.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		inspect.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		nilness.Analyzer,
		pkgfact.Analyzer,
		printf.Analyzer,
		reflectvaluecompare.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		slog.Analyzer,
		sortslice.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		tests.Analyzer,
		timeformat.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
		unusedwrite.Analyzer,
		usesgenerics.Analyzer,
	}
}

func staticcheckAnalyzers() []*analysis.Analyzer {
	var checksFromStaticcheck []*analysis.Analyzer

	// ST1...
	checksFromStaticcheck = append(checksFromStaticcheck, stylecheck.Analyzers[0].Analyzer)

	// S1...
	checksFromStaticcheck = append(checksFromStaticcheck, simple.Analyzers[0].Analyzer)

	// QF1...
	checksFromStaticcheck = append(checksFromStaticcheck, quickfix.Analyzers[0].Analyzer)

	// SA...
	for _, v := range staticcheck.Analyzers {
		checksFromStaticcheck = append(checksFromStaticcheck, v.Analyzer)
	}

	return checksFromStaticcheck
}

func otherPublicAnalyzers() []*analysis.Analyzer {
	out := make([]*analysis.Analyzer, 0, 2)

	out = append(out, rowserr.NewAnalyzer())
	out = append(out, bodyclose.Analyzer)

	return out
}

func osExitChecker() *analysis.Analyzer {
	return osexitchecker.OsExitAnalyzer
}
