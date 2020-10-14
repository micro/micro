// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Original source: github.com/micro/micro/v3/router/route_test.go

package router

import "testing"

func TestHash(t *testing.T) {
	route1 := Route{
		Service: "dest.svc",
		Gateway: "dest.gw",
		Network: "dest.network",
		Link:    "det.link",
		Metric:  10,
	}

	// make a copy
	route2 := route1

	route1Hash := route1.Hash()
	route2Hash := route2.Hash()

	// we should get the same hash
	if route1Hash != route2Hash {
		t.Errorf("identical routes result in different hashes")
	}
}
