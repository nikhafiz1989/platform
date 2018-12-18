package outputs

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/influxdata/flux"
	"github.com/influxdata/flux/execute"
	"github.com/influxdata/flux/plan"
	"github.com/influxdata/flux/semantic"
)

// NotifyAppendFileKind is the kind for the `notifyAppendFile` flux function
const NotifyAppendFileKind = "notifyAppendFile"

func init() {
	notifyAppendFileSignature := flux.FunctionSignature(
		map[string]semantic.PolyType{
			"path": semantic.String,
		},
		[]string{"path", flux.TablesParameter},
	)

	flux.RegisterFunctionWithSideEffect(NotifyAppendFileKind, createNotifyAppendFileOpSpec, notifyAppendFileSignature)
	flux.RegisterOpSpec(NotifyAppendFileKind, func() flux.OperationSpec { return &NotifyAppendFileOpSpec{} })
	plan.RegisterProcedureSpecWithSideEffect(NotifyAppendFileKind, newNotifyAppendFileProcedure, NotifyAppendFileKind)
	execute.RegisterTransformation(NotifyAppendFileKind, createNotifyAppendFileTransformation)
}

// NotifyAppendFileOpSpec is the flux.OperationSpec for notifyAppendFile.
type NotifyAppendFileOpSpec struct {
	Path string `json:"path"`
}

func (NotifyAppendFileOpSpec) Kind() flux.OperationKind {
	return NotifyAppendFileKind
}

var notifyAppendBaseDir string

func init() {
	dir := os.Getenv("INFLUXDB_NOTIFYDIR")
	if dir == "" {
		panic("Must set INFLUXDB_NOTIFYDIR env var to continue with notifyAppendFile spike")
	}

	if err := os.MkdirAll(dir, 0744); err != nil {
		panic(err)
	}

	notifyAppendBaseDir = dir
}

func createNotifyAppendFileOpSpec(args flux.Arguments, a *flux.Administration) (flux.OperationSpec, error) {
	if err := a.AddParentFromArgs(args); err != nil {
		return nil, err
	}

	rawPath, err := args.GetRequiredString("path")
	if err != nil {
		return nil, err
	}

	cleanPath := filepath.Clean(rawPath)
	return &NotifyAppendFileOpSpec{Path: filepath.Join(notifyAppendBaseDir, cleanPath)}, nil
}

func newNotifyAppendFileProcedure(qs flux.OperationSpec, a plan.Administration) (plan.ProcedureSpec, error) {
	spec, ok := qs.(*NotifyAppendFileOpSpec)
	if !ok && spec != nil {
		return nil, fmt.Errorf("invalid spec type %T", qs)
	}
	return &NotifyAppendFileProcedureSpec{Spec: spec}, nil
}

type NotifyAppendFileProcedureSpec struct {
	Spec *NotifyAppendFileOpSpec
}

func (o *NotifyAppendFileProcedureSpec) Copy() plan.ProcedureSpec {
	return &NotifyAppendFileProcedureSpec{
		Spec: &NotifyAppendFileOpSpec{
			Path: o.Spec.Path,
		},
	}
}

func (*NotifyAppendFileProcedureSpec) Kind() plan.ProcedureKind {
	return NotifyAppendFileKind
}

func createNotifyAppendFileTransformation(id execute.DatasetID, mode execute.AccumulationMode, spec plan.ProcedureSpec, a execute.Administration) (execute.Transformation, execute.Dataset, error) {
	s, ok := spec.(*NotifyAppendFileProcedureSpec)
	if !ok {
		return nil, nil, fmt.Errorf("invalid spec type %T", spec)
	}
	cache := execute.NewTableBuilderCache(a.Allocator())
	d := execute.NewDataset(id, mode, cache)

	t, err := NewNotifyAppendFileTransformation(d, cache, s)
	if err != nil {
		return nil, nil, err
	}
	return t, d, nil
}

type NotifyAppendFileTransformation struct {
	d     execute.Dataset
	cache execute.TableBuilderCache
	spec  *NotifyAppendFileProcedureSpec
}

func NewNotifyAppendFileTransformation(d execute.Dataset, cache execute.TableBuilderCache, spec *NotifyAppendFileProcedureSpec) (*NotifyAppendFileTransformation, error) {
	return &NotifyAppendFileTransformation{
		d:     d,
		cache: cache,
		spec:  spec,
	}, nil
}

func (t *NotifyAppendFileTransformation) Process(_ execute.DatasetID, tbl flux.Table) error {
	f, err := os.OpenFile(t.spec.Spec.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Println("NotifyAppendFile: appending to " + f.Name())

	fmt.Fprintf(f, "[%s] appending table with key %s\n", time.Now().Format(time.RFC3339Nano), tbl.Key())
	if err := tbl.Do(func(cr flux.ColReader) error {
		if _, err := fmt.Fprintf(f, "\t%s\n", cr.Key()); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return f.Close()
}

func (t *NotifyAppendFileTransformation) RetractTable(_ execute.DatasetID, key flux.GroupKey) error {
	return t.d.RetractTable(key)
}

func (t *NotifyAppendFileTransformation) Finish(_ execute.DatasetID, err error) {
	t.d.Finish(err)
}

func (t *NotifyAppendFileTransformation) UpdateProcessingTime(_ execute.DatasetID, pt execute.Time) error {
	return t.d.UpdateProcessingTime(pt)
}

func (t *NotifyAppendFileTransformation) UpdateWatermark(id execute.DatasetID, pt execute.Time) error {
	return t.d.UpdateWatermark(pt)
}
