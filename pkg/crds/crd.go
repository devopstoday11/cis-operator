package crds

import (
	"fmt"
	"io/ioutil"
	"strings"

	cisoperator "github.com/rancher/cis-operator/pkg/apis/cis.cattle.io/v1"
	"github.com/rancher/wrangler/pkg/crd"
	_ "github.com/rancher/wrangler/pkg/generated/controllers/apiextensions.k8s.io" //using init
	"github.com/rancher/wrangler/pkg/yaml"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func WriteCRD() error {
	for _, crdDef := range List() {
		bCrd, err := crdDef.ToCustomResourceDefinition()
		if err != nil {
			return err
		}
		yamlBytes, err := yaml.Export(&bCrd)
		if err != nil {
			return err
		}

		filename := fmt.Sprintf("./crds/%s.yaml", strings.ToLower(bCrd.Spec.Names.Kind))
		err = ioutil.WriteFile(filename, yamlBytes, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func List() []crd.CRD {
	return []crd.CRD{
		newCRD(&cisoperator.ClusterScan{}, func(c crd.CRD) crd.CRD {
			return c.
				WithColumn("ClusterScanProfile", ".status.lastRunScanProfileName").
				WithColumn("Total", ".status.summary.total").
				WithColumn("Pass", ".status.summary.pass").
				WithColumn("Fail", ".status.summary.fail").
				WithColumn("Skip", ".status.summary.skip").
				WithColumn("Not Applicable", ".status.summary.notApplicable").
				WithColumn("LastRunTimestamp", ".status.lastRunTimestamp")
		}),
		newCRD(&cisoperator.ClusterScanProfile{}, func(c crd.CRD) crd.CRD {
			return c.
				WithColumn("BenchmarkVersion", ".spec.benchmarkVersion")
		}),
		newCRD(&cisoperator.ClusterScanReport{}, func(c crd.CRD) crd.CRD {
			return c.
				WithColumn("LastRunTimestamp", ".spec.lastRunTimestamp").
				WithColumn("BenchmarkVersion", ".spec.benchmarkVersion")
		}),
		newCRD(&cisoperator.ClusterScanBenchmark{}, func(c crd.CRD) crd.CRD {
			return c.
				WithColumn("ClusterProvider", ".spec.clusterProvider").
				WithColumn("MinKubernetesVersion", ".spec.minKubernetesVersion").
				WithColumn("MaxKubernetesVersion", ".spec.maxKubernetesVersion")
		}),
	}
}

func newCRD(obj interface{}, customize func(crd.CRD) crd.CRD) crd.CRD {
	crd := crd.CRD{
		GVK: schema.GroupVersionKind{
			Group:   "cis.cattle.io",
			Version: "v1",
		},
		NonNamespace: true,
		Status:       true,
		SchemaObject: obj,
	}
	if customize != nil {
		crd = customize(crd)
	}
	return crd
}
