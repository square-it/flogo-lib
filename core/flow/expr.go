package flow

import "github.com/TIBCOSoftware/flogo-lib/core/data"

//todo rename
type LinkExprManager interface {

	// EvalLinkExpr evaluate the link expression
	EvalLinkExpr(link *Link, scope data.Scope) bool
}

func GetExpressionLinks(def *Definition) []*Link {

	var links []*Link
	getExpressionLinks(def.RootTask(), &links);

	return links;
}

func getExpressionLinks(task *Task, links *[]*Link) {

	for _, link := range  task.ChildLinks() {

		if link.Type() == LtExpression {
			*links = append(*links, link)
		}
	}

	for _, childTask := range task.ChildTasks() {
		getExpressionLinks(childTask, links)
	}
}

