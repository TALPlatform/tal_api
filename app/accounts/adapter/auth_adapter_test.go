package adapter

import (
	"testing"

	"github.com/TALPlatform/tal_api/db"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
	"github.com/stretchr/testify/require"
)

var adapterInstance AccountsAdapterInterface = NewAccountsAdapter()

func isDeepEqual(t *testing.T, item1 *talv1.NavigationBarItem, item2 *talv1.NavigationBarItem) bool {
	require.Equal(t, item1.Label, item2.Label)
	require.Equal(t, item1.Key, item2.Key)
	require.Equal(t, item1.Level, item2.Level)
	require.Equal(t, item1.NavigationBarItemId, item2.NavigationBarItemId)
	require.Equal(t, item1.ParentId, item2.ParentId)
	require.Equal(t, len(item1.Items), len(item2.Items))
	if len(item1.Items) > 0 {
		for i := range item1.Items {
			if !isDeepEqual(t, item1.Items[i], item2.Items[i]) {
				return false
			}
		}

	}
	return true
}
func executeTest(t *testing.T, req []db.UserNavigationBarFindRow, expected []*talv1.NavigationBarItem) {
	result, err := adapterInstance.UserNavigationBarFindGrpcFromSql(req, nil)
	require.NoError(t, err)
	require.Equal(t, len(expected), len(result))
	for i, acutalRecord := range result {
		expectedRecord := expected[i]
		require.True(t, isDeepEqual(t, acutalRecord, expectedRecord))
	}

}

type testCase struct {
	name     string
	req      []db.UserNavigationBarFindRow
	expected []*talv1.NavigationBarItem
}

func TestUserNavigationBarFindGrpcFromSql(t *testing.T) {
	testCases := []testCase{

		{
			name: "singleLevel",
			req: []db.UserNavigationBarFindRow{
				{
					NavigationBarItemID: 1,
					MenuKey:             "home",
					Label:               "Home",
					ParentID:            db.ToPgInt(0),
					Level:               1,
				},
			},
			expected: []*talv1.NavigationBarItem{
				{
					NavigationBarItemId: 1,
					ParentId:            0,
					Key:                 "home",
					Label:               "Home",
					Level:               1,
					Items:               nil,
				},
			},
		},
		{
			name: "twoLevels",
			req: []db.UserNavigationBarFindRow{
				{
					NavigationBarItemID: 1,
					MenuKey:             "home",
					Label:               "Home",
					ParentID:            db.ToPgInt(0),
					Level:               1,
				},
				{
					NavigationBarItemID: 2,
					MenuKey:             "dashboard",
					Label:               "Dashboard",
					ParentID:            db.ToPgInt(1),
					Level:               2,
				},
			},
			expected: []*talv1.NavigationBarItem{
				{
					NavigationBarItemId: 1,
					ParentId:            0,
					Key:                 "home",
					Label:               "Home",
					Level:               1,
					Items: []*talv1.NavigationBarItem{
						{
							NavigationBarItemId: 2,
							ParentId:            1,
							Key:                 "dashboard",
							Label:               "Dashboard",
							Level:               2,
						},
					},
				},
			},
		},
		{
			name: "threeLevels",
			req: []db.UserNavigationBarFindRow{
				{
					NavigationBarItemID: 1,
					MenuKey:             "home",
					Label:               "Home",
					ParentID:            db.ToPgInt(0),
					Level:               1,
				},
				{
					NavigationBarItemID: 2,
					MenuKey:             "dashboard",
					Label:               "Dashboard",
					ParentID:            db.ToPgInt(1),
					Level:               2,
				},
				{
					NavigationBarItemID: 3,
					MenuKey:             "analytics",
					Label:               "Analytics",
					ParentID:            db.ToPgInt(2),
					Level:               3,
				},
			},
			expected: []*talv1.NavigationBarItem{
				{
					NavigationBarItemId: 1,
					ParentId:            0,
					Key:                 "home",
					Label:               "Home",
					Level:               1,
					Items: []*talv1.NavigationBarItem{
						{
							NavigationBarItemId: 2,
							ParentId:            1,
							Key:                 "dashboard",
							Label:               "Dashboard",
							Level:               2,
							Items: []*talv1.NavigationBarItem{
								{
									NavigationBarItemId: 3,
									ParentId:            2,
									Key:                 "analytics",
									Label:               "Analytics",
									Level:               3,
									Items:               nil,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "fourLevels",
			req: []db.UserNavigationBarFindRow{
				{
					NavigationBarItemID: 1,
					MenuKey:             "home",
					Label:               "Home",
					ParentID:            db.ToPgInt(0),
					Level:               1,
				},
				{
					NavigationBarItemID: 2,
					MenuKey:             "dashboard",
					Label:               "Dashboard",
					ParentID:            db.ToPgInt(1),
					Level:               2,
				},
				{
					NavigationBarItemID: 3,
					MenuKey:             "analytics",
					Label:               "Analytics",
					ParentID:            db.ToPgInt(2),
					Level:               3,
				},
				{
					NavigationBarItemID: 4,
					MenuKey:             "reports",
					Label:               "Reports",
					ParentID:            db.ToPgInt(3),
					Level:               4,
				},
			},
			expected: []*talv1.NavigationBarItem{
				{
					NavigationBarItemId: 1,
					ParentId:            0,
					Key:                 "home",
					Label:               "Home",
					Level:               1,
					Items: []*talv1.NavigationBarItem{
						{
							NavigationBarItemId: 2,
							ParentId:            1,
							Key:                 "dashboard",
							Label:               "Dashboard",
							Level:               2,
							Items: []*talv1.NavigationBarItem{
								{
									NavigationBarItemId: 3,
									ParentId:            2,
									Key:                 "analytics",
									Label:               "Analytics",
									Level:               3,
									Items: []*talv1.NavigationBarItem{
										{
											NavigationBarItemId: 4,
											ParentId:            3,
											Key:                 "reports",
											Label:               "Reports",
											Level:               4,
											Items:               nil,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			executeTest(t, tc.req, tc.expected)
		})
	}

}
