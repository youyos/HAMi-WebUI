package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
	"time"
)

type ResourcePool struct {
	Id         int64     `db:"id"`
	PoolName   string    `db:"pool_name"`
	PoolType   int32     `db:"pool_type"`
	CreateTime time.Time `db:"create_time"`
	UpdateTime time.Time `db:"update_time"`
}

type Nodes struct {
	Id         int64     `db:"id"`
	NodeName   string    `db:"node_name"`
	NodeIp     string    `db:"node_ip"`
	CreateTime time.Time `db:"create_time"`
	UpdateTime time.Time `db:"update_time"`
}

type NodeInfo struct {
	Name string
	IP   string
}

func ExistsResourcePoolByPoolName(poolName string) bool {
	var count int
	err := db.QueryRow("SELECT count(1) FROM resource_pool WHERE pool_name = ?", poolName).Scan(&count)
	if err != nil {
		log.Infof("Query failed: %v", err)
		return false
	}

	return count > 0
}

func QueryResourcePoolById(poolId int64) (*ResourcePool, error) {
	var pool ResourcePool
	err := db.QueryRow("SELECT id, pool_name, pool_type, create_time, update_time FROM resource_pool WHERE id = ?", poolId).
		Scan(&pool.Id, &pool.PoolName, &pool.PoolType, &pool.CreateTime, &pool.UpdateTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Infof("No record found with id %d", poolId)
			return nil, nil
		}
		log.Infof("Query failed: %v", err)
		return nil, err
	}

	return &pool, nil
}

func QueryResourcePoolListAll() ([]*ResourcePool, error) {
	// 执行查询
	rows, err := db.Query("SELECT id, pool_name, pool_type, create_time, update_time FROM resource_pool order by create_time desc")
	if err != nil {
		log.Infof("Query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	// 存放结果的切片
	pools := make([]*ResourcePool, 0)

	// 遍历每一行
	for rows.Next() {
		var pool ResourcePool
		err := rows.Scan(&pool.Id, &pool.PoolName, &pool.PoolType, &pool.CreateTime, &pool.UpdateTime)
		if err != nil {
			log.Infof("Scan failed: %v", err)
			return nil, err
		}
		pools = append(pools, &pool)
	}

	// 检查 rows 是否遍历中出错
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pools, nil
}

func QueryNodesByPoolId(poolId int64) ([]*Nodes, error) {
	// 执行查询
	rows, err := db.Query("SELECT id, node_name, node_ip, create_time, update_time FROM nodes where pool_id = ?", poolId)
	if err != nil {
		log.Infof("Query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	// 存放结果的切片
	nodes := make([]*Nodes, 0)

	// 遍历每一行
	for rows.Next() {
		var node Nodes
		err := rows.Scan(&node.Id, &node.NodeName, &node.NodeIp, &node.CreateTime, &node.UpdateTime)
		if err != nil {
			log.Infof("Scan failed: %v", err)
			return nil, err
		}
		nodes = append(nodes, &node)
	}

	// 检查 rows 是否遍历中出错
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nodes, nil
}

func QueryDistinctNodes() ([]*Nodes, error) {
	// 执行查询
	rows, err := db.Query("select distinct nodes.node_name, nodes.node_ip from nodes")
	if err != nil {
		log.Infof("Query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	// 存放结果的切片
	nodes := make([]*Nodes, 0)

	// 遍历每一行
	for rows.Next() {
		var node Nodes
		err := rows.Scan(&node.NodeName, &node.NodeIp)
		if err != nil {
			log.Infof("Scan failed: %v", err)
			return nil, err
		}
		nodes = append(nodes, &node)
	}

	// 检查 rows 是否遍历中出错
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nodes, nil
}

func QueryResourceNamesByIp(nodeIp string) ([]string, error) {
	// 执行查询
	rows, err := db.Query("select pool_name from resource_pool where id in (select distinct pool_id from nodes where node_ip=?)", nodeIp)
	if err != nil {
		log.Infof("Query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	// 存放结果的切片
	resourcePoolNames := make([]string, 0)

	// 遍历每一行
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			log.Infof("Scan failed: %v", err)
			return nil, err
		}
		resourcePoolNames = append(resourcePoolNames, name)
	}

	// 检查 rows 是否遍历中出错
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return resourcePoolNames, nil
}

func QueryResourceNamesByNodeName(nodeName string) ([]string, error) {
	// 执行查询
	rows, err := db.Query("select pool_name from resource_pool where id in (select distinct pool_id from nodes where node_name=?)", nodeName)
	if err != nil {
		log.Infof("Query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	// 存放结果的切片
	resourcePoolNames := make([]string, 0)

	// 遍历每一行
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			log.Infof("Scan failed: %v", err)
			return nil, err
		}
		resourcePoolNames = append(resourcePoolNames, name)
	}

	// 检查 rows 是否遍历中出错
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return resourcePoolNames, nil
}

func InsertResourcePool(poolName string, poolType int32) (int64, error) {
	querySql := "INSERT INTO resource_pool(pool_name, pool_type) VALUES (?, ?)"

	result, err := db.Exec(querySql, poolName, poolType)
	if err != nil {
		log.Infof("Failed to insert record: %v", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Infof("Failed to get last insert ID: %v", err)
		return 0, err
	}

	return id, nil
}

func UpdateResourcePool(poolId int64, poolName string, poolType int32) (int64, error) {
	updateSql := "UPDATE resource_pool SET pool_name=?, pool_type=? where id=?"
	result, err := db.Exec(updateSql, poolName, poolType, poolId)
	if err != nil {
		log.Infof("Failed to update record: %v", err)
		return 0, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Infof("Failed to get rows affected: %v", err)
		return 0, err
	}

	return rows, nil
}

func InsertNodes(poolId int64, nodes []*NodeInfo) (int64, error) {
	valueStrings := make([]string, 0, len(nodes))
	valueArgs := make([]interface{}, 0, len(nodes)*3)

	for _, node := range nodes {
		valueStrings = append(valueStrings, "(?, ?, ?)")
		valueArgs = append(valueArgs, poolId, node.Name, node.IP)
	}

	insertSql := fmt.Sprintf("INSERT INTO nodes(pool_id, node_name, node_ip) VALUES %s",
		strings.Join(valueStrings, ","),
	)

	log.Info("InsertNodes: ", insertSql)
	result, err := db.Exec(insertSql, valueArgs...)
	if err != nil {
		log.Infof("Batch insert failed: %v", err)
		return 0, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Infof("Get rows affected failed: %v", err)
		return 0, err
	}

	return rows, nil
}

func DeleteResourcePoolById(poolId int64) (int64, error) {
	result, err := db.Exec("DELETE FROM resource_pool WHERE id = ?", poolId)
	if err != nil {
		return 0, fmt.Errorf("delete failed: %w", err)
	}

	// 返回影响的行数（0 表示未删除任何数据）
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("get rows affected failed: %w", err)
	}

	return rowsAffected, nil
}

func DeleteNodesByPoolId(poolId int64) (int64, error) {
	result, err := db.Exec("DELETE FROM nodes WHERE pool_id = ?", poolId)
	if err != nil {
		return 0, fmt.Errorf("delete failed: %w", err)
	}

	// 返回影响的行数（0 表示未删除任何数据）
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("get rows affected failed: %w", err)
	}

	return rowsAffected, nil
}

func DeleteNodeById(nodeId int64) (int64, error) {
	result, err := db.Exec("DELETE FROM nodes WHERE id = ?", nodeId)
	if err != nil {
		return 0, fmt.Errorf("delete failed: %w", err)
	}

	// 返回影响的行数（0 表示未删除任何数据）
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("get rows affected failed: %w", err)
	}

	return rowsAffected, nil
}
