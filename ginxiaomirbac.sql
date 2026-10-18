/*
 Navicat MySQL Data Transfer

 Source Server         : 192.168.0.6
 Source Server Type    : MySQL
 Source Server Version : 80020
 Source Host           : 192.168.0.6:3306
 Source Schema         : ginxiaomirbac

 Target Server Type    : MySQL
 Target Server Version : 80020
 File Encoding         : 65001
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for access
-- ----------------------------
DROP TABLE IF EXISTS `access`;
CREATE TABLE `access`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `module_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '0',
  `type` tinyint(1) NULL DEFAULT NULL,
  `action_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `url` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `module_id` int NULL DEFAULT NULL,
  `sort` int NULL DEFAULT NULL,
  `description` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `add_time` int NULL DEFAULT NULL,
  `status` int NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 118 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of access
-- ----------------------------
INSERT INTO `access` VALUES (52, '管理员管理', 1, '', '', 0, 105, '管理员管理', 0, 1);
INSERT INTO `access` VALUES (53, '角色管理', 1, '', '', 0, 100, '角色管理', 0, 1);
INSERT INTO `access` VALUES (54, '管理员管理', 2, '管理员列表', 'manager', 52, 100, '管理员列表', 0, 1);
INSERT INTO `access` VALUES (55, '管理员管理13', 2, '增加管理员', 'manager/add', 52, 101, '管理员列表', 0, 1);
INSERT INTO `access` VALUES (56, '管理员管理', 3, '编辑管理员', 'manager/edit', 52, 100, '编辑管理员', 0, 1);
INSERT INTO `access` VALUES (57, '管理员管理', 3, '删除管理员', 'manager/delete', 52, 100, '删除管理员', 0, 1);
INSERT INTO `access` VALUES (59, '角色管理', 2, '角色列表', 'role', 53, 100, '角色列表', 0, 1);
INSERT INTO `access` VALUES (60, '角色管理', 2, '增加角色', 'role/add', 53, 102, '增加角色', 0, 1);
INSERT INTO `access` VALUES (61, '角色管理', 3, '编辑角色', 'role/edit', 53, 100, '编辑角色', 0, 1);
INSERT INTO `access` VALUES (62, '角色管理', 3, '删除角色', 'role/delete', 53, 100, '删除角色', 0, 1);
INSERT INTO `access` VALUES (63, '权限管理', 1, '', '', 0, 100, '权限管理', 0, 1);
INSERT INTO `access` VALUES (64, '权限管理', 2, '权限列表', 'access', 63, 100, '', 0, 1);
INSERT INTO `access` VALUES (67, '权限管理', 2, '增加权限', 'access/add', 63, 100, '', 0, 1);
INSERT INTO `access` VALUES (69, '系统设置', 2, '轮播图列表', 'focus', 106, 101, '1111', 0, 1);
INSERT INTO `access` VALUES (70, '系统设置', 3, '增加轮播图', 'focus/add', 106, 100, '增加轮播图', 0, 1);
INSERT INTO `access` VALUES (71, '系统设置', 3, '编辑轮播图', 'focus/edit', 106, 100, '', 0, 1);
INSERT INTO `access` VALUES (75, '系统设置', 3, '删除轮播图', 'focus/delete', 106, 100, '', 0, 1);
INSERT INTO `access` VALUES (76, '管理员管理', 3, '执行增加管理员', 'manager/doAdd', 52, 100, '执行增加', 0, 1);
INSERT INTO `access` VALUES (77, '管理员管理', 3, '执行修改管理员', 'manager/doEdit', 52, 100, '执行修改', 0, 1);
INSERT INTO `access` VALUES (78, '角色管理', 3, '执行增加角色', 'role/doAdd', 53, 100, '执行增加', 0, 1);
INSERT INTO `access` VALUES (79, '角色管理', 3, '执行修改角色', 'role/doEdit', 53, 100, '执行修改', 0, 1);
INSERT INTO `access` VALUES (80, '角色管理', 3, '角色授权', 'role/auth', 53, 100, '', 0, 1);
INSERT INTO `access` VALUES (81, '角色管理', 3, '执行角色授权', 'role/doAuth', 53, 100, '执行授权', 0, 1);
INSERT INTO `access` VALUES (82, '权限管理', 3, '修改权限', 'access/edit', 63, 100, '执行修改', 0, 1);
INSERT INTO `access` VALUES (83, '权限管理', 3, '删除权限', 'access/delete', 63, 100, '', 0, 1);
INSERT INTO `access` VALUES (84, '权限管理', 3, '执行增加权限', 'access/doAdd', 63, 100, '', 0, 1);
INSERT INTO `access` VALUES (85, '权限管理', 3, '执行修改权限', 'access/doEdit', 63, 100, '执行修改\r\n', 0, 1);
INSERT INTO `access` VALUES (86, '系统设置', 3, '执行增加', 'focus/doAdd', 106, 100, '', 0, 1);
INSERT INTO `access` VALUES (87, '商品管理', 1, '', '', 0, 100, '', 0, 1);
INSERT INTO `access` VALUES (88, '商品管理', 2, '商品分类列表', 'goodsCate', 87, 100, '', 0, 1);
INSERT INTO `access` VALUES (89, '商品管理', 3, '增加商品分类', 'goodsCate/add', 87, 100, '', 0, 1);
INSERT INTO `access` VALUES (90, '商品管理', 3, '执行增加商品分类', 'goodsCate/doAdd', 87, 100, '', 0, 1);
INSERT INTO `access` VALUES (91, '商品管理', 3, '修改商品分类', 'goodsCate/edit', 87, 100, '', 0, 1);
INSERT INTO `access` VALUES (92, '商品管理', 3, '执行修改商品分类', 'goodsCate/doEdit', 87, 100, '', 0, 1);
INSERT INTO `access` VALUES (93, '商品管理', 3, '删除商品分类', 'goodsCate/delete', 87, 100, '', 0, 1);
INSERT INTO `access` VALUES (94, '商品管理', 2, '商品类型列表', 'goodsType', 87, 100, '', 0, 1);
INSERT INTO `access` VALUES (95, '商品管理', 3, '增加商品类型', 'goodsType/add', 87, 100, '', 0, 1);
INSERT INTO `access` VALUES (96, '商品管理', 3, '编辑商品类型', 'goodsType/edit', 87, 100, '', 0, 1);
INSERT INTO `access` VALUES (97, '商品管理', 3, '执行增加 商品类型', 'goodsType/doAdd', 87, 100, '', 0, 1);
INSERT INTO `access` VALUES (98, '商品管理', 3, '执行修改 商品类型', 'goodsType/doEdit', 87, 100, '', 0, 1);
INSERT INTO `access` VALUES (99, '商品管理', 3, '删除 商品类型', 'goodsType/delete', 87, 100, '', 0, 1);
INSERT INTO `access` VALUES (100, '商品管理', 2, '商品列表', 'goods', 87, 100, '商品列表', 0, 1);
INSERT INTO `access` VALUES (101, '商品管理', 3, '增加商品', 'goods/add', 87, 100, '', 0, 1);
INSERT INTO `access` VALUES (102, '商品管理', 3, '执行 增加商品', 'goods/doAdd', 87, 100, '', 0, 1);
INSERT INTO `access` VALUES (103, '商品管理', 3, '修改商品', 'goods/edit', 87, 100, '', 0, 1);
INSERT INTO `access` VALUES (104, '商品管理', 3, '执行 修改商品', 'goods/doEdit', 87, 100, '', 0, 1);
INSERT INTO `access` VALUES (105, '商品管理', 3, '删除商品', 'goods/delete', 87, 100, '', 0, 1);
INSERT INTO `access` VALUES (106, '系统设置', 1, '', '', 0, 100, '', 0, 1);
INSERT INTO `access` VALUES (107, '系统设置', 2, '导航管理', 'nav', 106, 100, '', 0, 1);
INSERT INTO `access` VALUES (108, '系统设置', 3, '增加导航', 'nav/add', 106, 100, '', 0, 1);
INSERT INTO `access` VALUES (109, '系统设置', 3, '编辑导航', 'nav/edit', 106, 100, '', 0, 1);
INSERT INTO `access` VALUES (110, '系统设置', 3, '删除导航', 'nav/delete', 106, 100, '', 0, 1);
INSERT INTO `access` VALUES (111, '系统设置', 3, '执行增加', 'nav/doAdd', 106, 100, '', 0, 1);
INSERT INTO `access` VALUES (112, '系统设置', 3, '执行修改', 'nav/doEdit', 106, 100, '', 0, 1);
INSERT INTO `access` VALUES (113, '系统设置', 2, '商店设置', 'setting', 106, 100, '', 0, 1);

-- ----------------------------
-- Table structure for manager
-- ----------------------------
DROP TABLE IF EXISTS `manager`;
CREATE TABLE `manager`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `password` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `mobile` varchar(11) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `email` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `status` tinyint(1) NULL DEFAULT NULL,
  `role_id` int NULL DEFAULT NULL,
  `add_time` int NULL DEFAULT NULL,
  `is_super` tinyint(1) NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `role_id`(`role_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 16 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of manager
-- ----------------------------
INSERT INTO `manager` VALUES (1, 'admin', 'e10adc3949ba59abbe56e057f20f883e', '152016122', '51886112214@qq.com', 1, 9, 0, 1);
INSERT INTO `manager` VALUES (2, 'zhangsan', 'fcea920f7412b5da7be0cf42b8c93759', '1520111122', '342338691122@qq.com', 1, 14, 1581661532, 0);
INSERT INTO `manager` VALUES (6, 'lisi', 'e10adc3949ba59abbe56e057f20f883e', '1520171122', '11114292@qq.com', 1, 16, 1631156378, 0);

-- ----------------------------
-- Table structure for role
-- ----------------------------
DROP TABLE IF EXISTS `role`;
CREATE TABLE `role`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `status` tinyint(1) NULL DEFAULT NULL,
  `add_time` int NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 20 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of role
-- ----------------------------
INSERT INTO `role` VALUES (9, '超级管理员', '我是一个超级管理员', 1, 1631072961);
INSERT INTO `role` VALUES (14, '软件部门', '软件部门', 1, 1631075350);
INSERT INTO `role` VALUES (16, '销售部门3', '销售部门', 1, 1631589828);

-- ----------------------------
-- Table structure for role_access
-- ----------------------------
DROP TABLE IF EXISTS `role_access`;
CREATE TABLE `role_access`  (
  `role_id` int NOT NULL,
  `access_id` int NOT NULL,
  INDEX `role_id`(`role_id`) USING BTREE,
  INDEX `access_id`(`access_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of role_access
-- ----------------------------
INSERT INTO `role_access` VALUES (9, 52);
INSERT INTO `role_access` VALUES (9, 54);
INSERT INTO `role_access` VALUES (9, 55);
INSERT INTO `role_access` VALUES (9, 53);
INSERT INTO `role_access` VALUES (9, 59);
INSERT INTO `role_access` VALUES (9, 60);
INSERT INTO `role_access` VALUES (9, 61);
INSERT INTO `role_access` VALUES (9, 62);
INSERT INTO `role_access` VALUES (9, 63);
INSERT INTO `role_access` VALUES (9, 64);
INSERT INTO `role_access` VALUES (9, 67);
INSERT INTO `role_access` VALUES (9, 82);
INSERT INTO `role_access` VALUES (9, 83);
INSERT INTO `role_access` VALUES (9, 84);
INSERT INTO `role_access` VALUES (9, 85);
INSERT INTO `role_access` VALUES (9, 70);
INSERT INTO `role_access` VALUES (9, 71);
INSERT INTO `role_access` VALUES (16, 53);
INSERT INTO `role_access` VALUES (16, 59);
INSERT INTO `role_access` VALUES (16, 60);
INSERT INTO `role_access` VALUES (16, 61);
INSERT INTO `role_access` VALUES (16, 62);
INSERT INTO `role_access` VALUES (16, 78);
INSERT INTO `role_access` VALUES (16, 79);
INSERT INTO `role_access` VALUES (16, 80);
INSERT INTO `role_access` VALUES (16, 81);
INSERT INTO `role_access` VALUES (16, 63);
INSERT INTO `role_access` VALUES (16, 64);
INSERT INTO `role_access` VALUES (16, 67);
INSERT INTO `role_access` VALUES (16, 82);
INSERT INTO `role_access` VALUES (16, 83);
INSERT INTO `role_access` VALUES (16, 84);
INSERT INTO `role_access` VALUES (16, 85);
INSERT INTO `role_access` VALUES (16, 87);
INSERT INTO `role_access` VALUES (16, 88);
INSERT INTO `role_access` VALUES (16, 89);
INSERT INTO `role_access` VALUES (14, 52);
INSERT INTO `role_access` VALUES (14, 55);
INSERT INTO `role_access` VALUES (14, 54);
INSERT INTO `role_access` VALUES (14, 56);
INSERT INTO `role_access` VALUES (14, 57);
INSERT INTO `role_access` VALUES (14, 76);
INSERT INTO `role_access` VALUES (14, 77);
INSERT INTO `role_access` VALUES (14, 53);
INSERT INTO `role_access` VALUES (14, 60);
INSERT INTO `role_access` VALUES (14, 59);
INSERT INTO `role_access` VALUES (14, 61);
INSERT INTO `role_access` VALUES (14, 62);
INSERT INTO `role_access` VALUES (14, 78);
INSERT INTO `role_access` VALUES (14, 79);
INSERT INTO `role_access` VALUES (14, 80);
INSERT INTO `role_access` VALUES (14, 81);
INSERT INTO `role_access` VALUES (14, 63);
INSERT INTO `role_access` VALUES (14, 64);
INSERT INTO `role_access` VALUES (14, 67);
INSERT INTO `role_access` VALUES (14, 82);
INSERT INTO `role_access` VALUES (14, 83);
INSERT INTO `role_access` VALUES (14, 84);
INSERT INTO `role_access` VALUES (14, 85);

SET FOREIGN_KEY_CHECKS = 1;
