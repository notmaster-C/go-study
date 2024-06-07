package utils

import (
	"fmt"
	"go-study/db"
	"io/ioutil"
	"path/filepath"
)

func DbflFile() {
	rows, err := db.GetDb("mysql").Table("T_GXZS_ZSXXB_X").Where("XQ=?", "校区").Rows()
	if err != nil {
		fmt.Println("数据库查询错误")
	}
	// 根据查询结果遍历并移动文件
	for rows.Next() {
		var sfzh string
		err := rows.Scan(&sfzh)
		if err != nil {
			fmt.Println("数据读取错误:", err)
			continue
		}

		// 源文件路径（根据实际情况修改）
		sourcePath := "C:\\Users\\Administrator\\Pictures\\zp2\\file\\" + sfzh + ".jpg"

		// 目标文件夹路径
		targetPath := "C:\\Users\\Administrator\\Pictures\\pd\\"

		// 创建分类文件夹
		//err = os.MkdirAll(targetPath, 0755)
		//if err != nil {
		//	fmt.Println("创建文件夹失败:", err)
		//	continue
		//}

		// 移动文件
		err = moveFile(sourcePath, filepath.Join(targetPath, sfzh+".jpg"))
		if err != nil {
			fmt.Println("移动文件失败:", err)
			continue
		}

		fmt.Println("文件移动成功:", sfzh)
	}
}

// 移动文件
func moveFile(sourcePath, targetPath string) error {
	fileBytes, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(targetPath, fileBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
