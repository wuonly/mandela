package dao

import (
	"fmt"
	"math/big"
	// "strconv"
)

type Bucket struct {
	dimension int     //1为左节点，0为右节点
	Bucket    []Node  //保存的节点
	MaxSize   int     //k桶大小
	level     int     //深度 根节点深度为512,子节点深度最小为0
	left      *Bucket //左边1
	right     *Bucket //右边0
	parent    *Bucket //父节点
}

func (this *Bucket) getRootIdLoop() *Bucket {
	if this.parent == nil {
		return this
	} else {
		return this.getRootIdLoop()
	}
}

//得到根节点的id
func (this *Bucket) GetRootId() string {
	return this.getRootIdLoop().Bucket[0].NodeId
}

//插入一个节点，将一个节点保存到树中的指定位置
func (this *Bucket) insertLoop(nodeSelf, nodeInsert *big.Int, node Node) *Bucket {
	//-----------------------------
	//若深度为0，则返回当前节点
	//-----------------------------
	if this.level == 0 {
		if this.Bucket == nil {
			this.Bucket = []Node{node}
		} else {
			this.Bucket = append(this.Bucket, node)
		}
		return this
	}
	//-----------------------------
	//如果最高位不相同，则返回本节点
	//-----------------------------
	if (nodeSelf.BitLen() == this.level || nodeInsert.BitLen() == this.level) && nodeSelf.BitLen() != nodeInsert.BitLen() {
		if this.Bucket == nil {
			this.Bucket = []Node{node}
		} else {
			this.Bucket = append(this.Bucket, node)
		}
		return this
	}
	//-----------------------------
	//截取最高位
	//-----------------------------
	tempInt := big.NewInt(1)
	tempInt = tempInt.Lsh(tempInt, uint(this.level-1))
	selfInt := new(big.Int).AndNot(nodeSelf, tempInt)
	insertInt := new(big.Int).AndNot(nodeInsert, tempInt)
	//-----------------------------
	//按最高位是1还是0添加一个左节点或者右节点
	//-----------------------------
	insertBucket := &Bucket{parent: this, level: this.level - 1}
	if nodeInsert.BitLen() == this.level {
		//最高位是1，添加一个左节点
		if this.left == nil {
			insertBucket.insertLoop(selfInt, insertInt, node)
			insertBucket.dimension = 1
			this.left = insertBucket
		} else {
			this.left.insertLoop(selfInt, insertInt, node)
		}
	} else {
		//最高位是0，添加一个右节点
		if this.right == nil {
			insertBucket.insertLoop(selfInt, insertInt, node)
			insertBucket.dimension = 0
			this.right = insertBucket
		} else {
			this.right.insertLoop(selfInt, insertInt, node)
		}
	}
	return this
}

//插入一个节点
//@node 待插入的节点
func (this *Bucket) InsertNode(node Node) {
	selfInt, _ := new(big.Int).SetString(this.GetRootId(), 10)
	insertInt, _ := new(big.Int).SetString(node.NodeId, 10)
	fmt.Println(node)
	fmt.Println(insertInt.String())
	this.insertLoop(selfInt, insertInt, node)
}

func (this *Bucket) findBucketLoop(nodeInt *big.Int, targetLevel int) *Bucket {
	//---------------------------------
	//  深度为0,则找到最小深度的子节点
	//---------------------------------
	if this.level == 0 || this.level == targetLevel {
		return this
	}
	//---------------------------------
	//  没有子节点了
	//---------------------------------
	if this.left == nil && this.right == nil {
		return this
	}
	//---------------------------------
	//  截取最高位
	//---------------------------------
	tempInt := big.NewInt(1)
	tempInt = tempInt.Lsh(tempInt, uint(this.level-1))
	findInt := new(big.Int).AndNot(nodeInt, tempInt)
	//判断最高位是1还是0
	if nodeInt.BitLen() == this.level {
		//最高位是1
		left := this.left.findBucketLoop(findInt, targetLevel)
		return left
	} else {
		//最高位是0
		right := this.right.findBucketLoop(findInt, targetLevel)
		return right
	}
}

//从子节点中找到本k桶的最近k桶
func (this *Bucket) findRecentChildBucketLoop() *Bucket {
	if len(this.Bucket) != 0 {
		return this
	}
	if this.left != nil {
		leftBucket := this.left.findRecentChildBucketLoop()
		if leftBucket != nil {
			return leftBucket
		}
	}
	if this.right != nil {
		rightBucket := this.right.findRecentChildBucketLoop()
		if rightBucket != nil {
			return rightBucket
		}
	}
	return nil
}

//从父节点中找到本k桶最近的k桶
func (this *Bucket) findRecentParentBucketLoop() *Bucket {
	if len(this.Bucket) != 0 {
		return this
	}
	//本节点为右节点
	if this.dimension == 0 {
		leftParentBucket := this.parent.left.findRecentChildBucketLoop()
		if leftParentBucket != nil {
			return leftParentBucket
		}
	} else {
		rightParentBucket := this.parent.right.findRecentChildBucketLoop()
		if rightParentBucket != nil {
			return rightParentBucket
		}
	}
	parentBucket := this.findRecentParentBucketLoop()
	return parentBucket
}

//找到某个节点最近的节点
//@nodeId  节点id的10进制字符串
func (this *Bucket) FindRecentNode(nodeId string) Node {
	findInt, _ := new(big.Int).SetString(nodeId, 10)
	rootInt, _ := new(big.Int).SetString(this.Bucket[0].NodeId, 10)
	targetLevelInt := new(big.Int).Xor(rootInt, findInt)

	bucket := this.findBucketLoop(findInt, targetLevelInt.BitLen())
	//这个k桶为空
	if len(bucket.Bucket) == 0 {
		//找子节点中的k桶
		childBucket := bucket.findRecentChildBucketLoop()
		if childBucket != nil {
			bucket = childBucket
		} else {
			bucket = bucket.findRecentParentBucketLoop()
		}
	}

	for _, node := range bucket.Bucket {
		if node.NodeId == nodeId {
			return node
		}
	}
	//如果是父节点k桶
	if bucket.Bucket[0].NodeId == this.GetRootId() {
		return Node{}
	}
	return bucket.Bucket[0]
}

//得到一个网络的所有id
//@pre              节点id的网络id，比如：10010010100000000000000000000000
//@supplementary    网络id后面0的个数
//@count            生成多少个id
func (this *Bucket) getBucketNodeIds(pre *big.Int, supplementary int, count int) []string {
	//最大范围
	maxRange := new(big.Int).Lsh(big.NewInt(1), uint(supplementary))
	//最大范围除以节点个数，得到每个节点的步长
	pacesInt := new(big.Int).Div(maxRange, big.NewInt(int64(count)))

	ids := []*big.Int{}
	ids = append(ids, pre)
	for i := 0; i < int(count-1); i++ {
		value := new(big.Int).Add(ids[len(ids)-1], pacesInt)
		ids = append(ids, value)
	}
	idStrs := []string{}
	for _, id := range ids {
		idStrs = append(idStrs, id.String())
	}
	return idStrs
}

//得到构建所有k桶所需要的每个节点id
func (this *Bucket) GetAllBucketNodeId() []string {
	rootInt, _ := new(big.Int).SetString(this.GetRootId(), 10)
	ids := []string{}
	for i := 0; i < this.level; i++ {
		//---------------------------------
		//最后一位取反
		//---------------------------------
		reverseRootInt := new(big.Int).Xor(rootInt, new(big.Int).Lsh(big.NewInt(1), uint(i)))
		countId := i
		//控制一个k桶最多装32个节点信息
		if i > 4 {
			countId = int(new(big.Int).Lsh(big.NewInt(1), uint(4)).Int64())
		} else {
			countId = int(new(big.Int).Lsh(big.NewInt(1), uint(i)).Int64())
		}
		resultIds := this.getBucketNodeIds(reverseRootInt, i, countId)
		for _, id := range resultIds {
			ids = append(ids, id)
		}
		//---------------------------------
		//截取最后一位
		//---------------------------------
		tempInt := new(big.Int).Rsh(big.NewInt(1), uint(i))
		rootInt = new(big.Int).AndNot(rootInt, tempInt)
	}
	return ids
}

func NewBucket(node Node, maxLevel int) *Bucket {
	Bucket := Bucket{Bucket: []Node{node}, level: maxLevel + 1}
	return &Bucket
}

//保存节点的id
//ip地址
//不同协议的端口
type Node struct {
	NodeId               string         //节点id的10进制字符串
	Addr                 string         //外网ip地址
	PortFromScheme       map[string]int `["TCP":1990,"UDP":1991]` //外网映射到端口
	LastContactTimestamp bool           //最后联系的时间戳
	BeingPinged          bool           //节点是否存在
}
