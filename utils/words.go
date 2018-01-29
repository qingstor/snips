// +-------------------------------------------------------------------------
// | Copyright (C) 2016 Yunify, Inc.
// +-------------------------------------------------------------------------
// | Licensed under the Apache License, Version 2.0 (the "License");
// | you may not use this work except in compliance with the License.
// | You may obtain a copy of the License in the LICENSE file, or at:
// |
// | http://www.apache.org/licenses/LICENSE-2.0
// |
// | Unless required by applicable law or agreed to in writing, software
// | distributed under the License is distributed on an "AS IS" BASIS,
// | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// | See the License for the specific language governing permissions and
// | limitations under the License.
// +-------------------------------------------------------------------------

package utils

import (
	"github.com/imdario/mergo"
)

var capitalizedToCapitalizedWordMap = map[string]string{
	"Dns":     "DNS",
	"Dyn":     "DYN",
	"Eip":     "EIP",
	"Keypair": "KeyPair",
	"Vxnet":   "VxNet",
}

var lowercaseToLowercaseWordMap = map[string]string{
	"lastest": "latest", // Fix typo
}

var lowercaseToCapitalizedWordMap = map[string]string{
	"acl":           "ACL",
	"cors":          "CORS",
	"cpu":           "CPU",
	"datadir":       "DataDir",
	"dhcp":          "DHCP",
	"dns":           "DNS",
	"dyn":           "DYN",
	"eip":           "EIP",
	"eips":          "EIPs",
	"http":          "HTTP",
	"icp":           "ICP",
	"id":            "ID",
	"ids":           "IDs",
	"innodb":        "InnoDB",
	"io":            "IO",
	"ip":            "IP",
	"ips":           "IPs",
	"ipset":         "IPSet",
	"ipsets":        "IPSets",
	"keypair":       "KeyPair",
	"keypairs":      "KeyPairs",
	"lastest":       "Latest", // Fix typo
	"loadbalancer":  "LoadBalancer",
	"loadbalancers": "LoadBalancers",
	"md5":           "MD5",
	"newsid":        "NewSID",
	"nic":           "NIC",
	"os":            "OS",
	"opt":           "OPT",
	"qingstor":      "QingStor",
	"qingcloud":     "QingCloud",
	"qs":            "QS",
	"rdb":           "RDB",
	"rdbs":          "RDBs",
	"sha1":          "SHA1",
	"sha256":        "SHA256",
	"sql":           "SQL",
	"tmp":           "TMP",
	"tmpdir":        "TMPDir",
	"topslave":      "TopSlave",
	"trx":           "TRX",
	"ui":            "UI",
	"uri":           "URI",
	"url":           "URL",
	"usb":           "USB",
	"uuid":          "UUID",
	"vcpus":         "VCPUs",
	"vxnet":         "VxNet",
	"vxnets":        "VxNets",
}

var abbreviateWordMap = []string{
	"ACL",
	"CORS",
	"CPU",
	"DHCP",
	"DNS",
	"DYN",
	"EIP",
	"ETag",
	"IaaS",
	"ICP",
	"ID",
	"IO",
	"IP",
	"IPSets",
	"MD5",
	"NIC",
	"OAuth",
	"OS",
	"OPT",
	"QingStor",
	"QingCloud",
	"QS",
	"RDB",
	"SQL",
	"SSO",
	"TMP",
	"TMPDir",
	"TRX",
	"UI",
	"URL",
	"UUID",
	"VCPUs",
	"VxNet",
}

// MergeCapitalizedToCapitalizedWordMap will merge capitalizedToCapitalizedWordMap.
func MergeCapitalizedToCapitalizedWordMap(m map[string]string) {
	err := mergo.MergeWithOverwrite(&capitalizedToCapitalizedWordMap, m)
	CheckErrorForExit(err, 1)
}

// MergeLowercaseToLowercaseWordMap will merge lowercaseToLowercaseWordMap.
func MergeLowercaseToLowercaseWordMap(m map[string]string) {
	err := mergo.MergeWithOverwrite(&lowercaseToLowercaseWordMap, m)
	CheckErrorForExit(err, 1)
}

// MergeLowercaseToCapitalizedWordMap will merge lowercaseToCapitalizedWordMap.
func MergeLowercaseToCapitalizedWordMap(m map[string]string) {
	err := mergo.MergeWithOverwrite(&lowercaseToCapitalizedWordMap, m)
	CheckErrorForExit(err, 1)
}

// MergeAbbreviateWordMap will merge abbreviateWordMap.
func MergeAbbreviateWordMap(m []string) {
	t := map[string]struct{}{}

	for _, v := range abbreviateWordMap {
		t[v] = struct{}{}
	}
	for _, v := range m {
		t[v] = struct{}{}
	}

	abbreviateWordMap = make([]string, len(t))
	i := 0
	for k := range t {
		abbreviateWordMap[i] = k
		i++
	}
}
