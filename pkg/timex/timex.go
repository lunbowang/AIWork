/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package timex

import "time"

func Format(date int64) string {
	return time.Unix(date, 0).Format("2006-01-02")
}
