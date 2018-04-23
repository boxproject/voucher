// Copyright 2017. box.la authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
pragma solidity ^0.4.10;

contract Oracle {
    struct Node {
        // 节点签名人
        address signer;
        // 是否被启用，可被停止
        bool enabled;
    }

    // 这个帐号初始时可以在一个节点上生成，生成后要在别的节点管理时，需要把该帐号挪过去。
    // 需要开发管理节点的接口来控制本合约
    address public boss;

    mapping (address => uint) nodeId;

    Node[] nodes;

    modifier onlyBoss {
        require(boss == msg.sender);
        _;
    }

    function Oracle() public {
        boss = msg.sender;
    }

    // 添加节点系统授权人，这个角色负责打包数据
    function addSigner(address signer) onlyBoss public {
        // init the array if empty
        if (nodes.length == 0) {
            nodes.length = 1;
            nodeId[0] = 0;
            nodes[0] = Node({signer:0, enabled: false});
        }

        uint id = nodeId[signer];
        if (id == 0) {
            id = nodes.length;
            nodeId[signer] = nodes.length;
            nodes.length++;
        }

        Node storage node = nodes[id];
        node.signer = signer;
        node.enabled = true;
        nodeId[signer] = id;
    }

    function disableSigner(address signer) onlyBoss public returns(bool) {
        uint id = nodeId[signer];
        if (id != 0) {
            nodes[id].enabled = false;
            return true;
        }

        return false;
    }

    function count() constant public returns(uint) {
        return nodes.length;
    }

    function indexOf(uint idx) constant public returns(address, bool) {
        address addr = nodes[idx].signer;
        bool enabled = nodes[idx].enabled;
        return (addr, enabled);
    }

    function totalEnabledNodes() constant public returns (uint) {
        uint len = nodes.length;
        uint i;
        for (i = 1; i < nodes.length; i++) {
            if (!nodes[i].enabled) {
                len -= 1;
            }
        }
        return len - 1;
    }

    function isSigner(address signer) constant public returns(bool) {
        return (nodeId[signer] != 0 && nodes[nodeId[signer]].enabled);
    }
}