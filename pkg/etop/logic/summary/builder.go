package summary

import (
	"errors"
	"fmt"

	"o.o/backend/pkg/common/sql/sq"
)

type SummaryQueryBuilder struct {
	Table        string
	Items        []*Subject
	ScanArgs     []interface{}
	SummaryJoins []*SummaryJoin
	Pred         Predicator
}

func NewSummaryQueryBuilder(table string) *SummaryQueryBuilder {
	return &SummaryQueryBuilder{
		Table: table,
	}
}

func NewSummaryQueryWithJoinsBuilder(table string, joins []*SummaryJoin) *SummaryQueryBuilder {
	return &SummaryQueryBuilder{
		Table:        table,
		SummaryJoins: joins,
	}
}

func (q *SummaryQueryBuilder) AddCell(s *Subject, scanArg interface{}) {
	q.Items = append(q.Items, s)
	q.ScanArgs = append(q.ScanArgs, scanArg)
}

func (q *SummaryQueryBuilder) WriteSQLTo(w sq.SQLWriter) error {
	if len(q.Items) == 0 {
		return errors.New("no item")
	}

	w.WriteRawString("SELECT \n")
	for _, item := range q.Items {
		if err := item.WriteSQLTo(w); err != nil {
			return err
		}
		w.WriteRawString(",\n")
	}
	w.TrimLast(2)
	w.WriteRawString("\nFROM ")
	w.WriteName(q.Table)

	if len(q.SummaryJoins) != 0 {
		w.WriteRawString("\n")
		for _, summaryJoin := range q.SummaryJoins {
			if err := summaryJoin.WriteSQLTo(w); err != nil {
				return err
			}
		}
	}
	if q.Pred != nil {
		w.WriteRawString(" WHERE ")
		if err := q.Pred.WriteSQLTo(w); err != nil {
			return err
		}
	}
	return nil
}

type Table struct {
	Label string
	Tags  []string
	NRow  int
	NCol  int

	Rows []Subject
	Cols []Predicator
	Data []Cell
}

func NewTable(row, col int, label string, tags ...string) *Table {
	if row <= 0 || col <= 0 {
		panic(fmt.Sprintf("invalid dimension %vx%v", row, col))
	}
	return &Table{
		Label: label,
		Tags:  tags,
		NRow:  row,
		NCol:  col,
		Data:  make([]Cell, row*col),
	}
}

func (t *Table) Clone(lable string, tags ...string) *Table {
	nt := NewTable(len(t.Rows), len(t.Cols), lable, tags...)
	nt.Rows = make([]Subject, len(t.Rows))
	nt.Cols = make([]Predicator, len(t.Cols))
	nt.Data = make([]Cell, len(t.Data))
	copy(nt.Rows, t.Rows)
	copy(nt.Cols, t.Cols)
	copy(nt.Data, t.Data)
	return nt
}

func (t *Table) Cell(row, col int) *Cell {
	return &t.Data[row*t.NCol+col]
}

func (t *Table) SetCell(row, col int, subject Subject) *Cell {
	cell := t.Cell(row, col)
	cell.Subject = subject
	return cell
}

type Cell struct {
	Subject Subject
	Value   int64
}

type Subject struct {
	Label string
	Spec  string
	Expr  string
	Unit  string
	Ident int

	Pred Predicator
}

func NewSubject(label string, unit string, spec string, expr string, pred Predicator) Subject {
	return Subject{
		Label: label,
		Spec:  spec,
		Expr:  expr,
		Unit:  unit,
		Pred:  pred,
	}
}

func (s Subject) Combine(label string, pred ...Predicator) Subject {
	var p Predicator
	if len(pred) != 0 {
		p = Compose("", pred...)
	}
	ns := Subject{
		Label: label,
		Spec:  s.Spec,
		Expr:  s.Expr,
		Unit:  s.Unit,
		Pred:  Compose("", s.Pred, p),
	}
	return ns
}

func (s Subject) WithIdent(ident int) Subject {
	s.Ident = ident // clone
	return s
}

func (s *Subject) GetLabel() string {
	return s.Label
}

func (s *Subject) GetSpec() string {
	if s.Pred == nil {
		return s.Spec + ":"
	}
	return s.Spec + ":" + s.Pred.GetSpec()
}

func (s *Subject) WriteSQLTo(w sq.SQLWriter) error {
	w.WriteRawString("    -- ")
	w.WriteRawString(s.GetSpec())
	w.WriteRawString("\n    ")
	w.WriteRawString(s.Expr)
	if s.Pred == nil {
		return nil
	}
	w.WriteRawString(" FILTER (WHERE ")
	err := s.Pred.WriteSQLTo(w)
	w.WriteRawString(")")
	return err
}

type Predicator interface {
	GetLabel() string
	GetSpec() string
	Clone(label string) Predicator
	sq.WriterTo
}

type Predicate struct {
	Label string
	Spec  string
	Expr  sq.WriterTo
}

func (p Predicate) GetLabel() string {
	return p.Label
}

func (p Predicate) GetSpec() string {
	return p.Spec
}

func (p Predicate) Clone(label string) Predicator {
	p.Label = label
	return p
}

func (p Predicate) WriteSQLTo(w sq.SQLWriter) error {
	return p.Expr.WriteSQLTo(w)
}

type ComposedPredicate struct {
	Label string
	Preds []Predicate
}

func (p ComposedPredicate) GetLabel() string {
	return p.Label
}

func (p ComposedPredicate) GetSpec() string {
	var b []byte
	for _, pred := range p.Preds {
		b = append(b, pred.GetSpec()...)
		b = append(b, ',')
	}
	return string(b[:len(b)-1])
}

func (p ComposedPredicate) Clone(label string) Predicator {
	preds := make([]Predicate, len(p.Preds))
	copy(preds, p.Preds)
	return ComposedPredicate{
		Label: label,
		Preds: preds,
	}
}

func (p ComposedPredicate) WriteSQLTo(w sq.SQLWriter) error {
	w.WriteRawString("(")
	for _, pred := range p.Preds {
		if err := pred.WriteSQLTo(w); err != nil {
			return err
		}
		w.WriteRawString(") AND (")
	}
	w.TrimLast(7)
	w.WriteRawString(")")
	return nil
}

type SummaryJoin struct {
	Table string
	Preds []string
}

func (p *SummaryJoin) WriteSQLTo(w sq.SQLWriter) error {
	w.WriteRawString("JOIN ")
	w.WriteName(p.Table)

	if len(p.Preds) != 0 {
		w.WriteRawString(" ON ")
		for i, pred := range p.Preds {
			w.WriteRawString(pred)
			if i != len(p.Preds)-1 {
				w.WriteRawString(" AND ")
			}
		}
	}
	w.WriteRawString("\n")
	return nil
}

func Compose(label string, preds ...Predicator) Predicator {
	switch len(preds) {
	case 0:
		return nil
	case 1:
		return preds[0].Clone(label)
	}
	var items []Predicate
	count := 0
	for _, p := range preds {
		if p == nil {
			continue
		}
		count++
		switch p := p.(type) {
		case Predicate:
			items = append(items, p)
		case ComposedPredicate:
			items = append(items, p.Preds...)
		default:
			panic("unexpected")
		}
		if label == "" {
			label = p.GetLabel()
		}
	}
	if count == 0 {
		return nil
	}
	return ComposedPredicate{
		Preds: items,
		Label: label,
	}
}
