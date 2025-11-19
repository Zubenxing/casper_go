## 目标
- 仅改动“对比分析”节点，保证两个月数据的对比按 `Resource` 唯一键执行。
- 以 `CostUSD` 为唯一的成本来源，针对同一 `Resource` 出现多行时先合计，再做差值对比。
- 保持导出的 CSV 字段不变：`Resource, SubscriptionName, ResourceGroup, CurrentCostUSD, LastCostUSD, DeltaUSD, Remark`。

## 方案 A：替换对比节点代码（推荐）
- 输入读取：显式读取两路输入，必要时启用回退按 `pairedItem.input` 分组，避免连接丢失时读错源。
- 键与聚合：
  - `key(row) = trim(row.Resource)`（满足你“Resource 唯一”的模型）
  - 分别对本月、上月聚合：对相同 `Resource` 的所有行将 `CostUSD` 合计为一个值（保留最近一行用于填充其他列）。
- 比较与标记：两边都有则计算 `DeltaUSD=Current-Last`，并标记 `Remark=UP/DOWN/SAME`；只在本月为 `NEW`，只在上月为 `REMOVED`。
- 代码（直接粘贴至“对比分析”节点，覆盖现有内容）：

```javascript
let cur = $input.all(0).map(i => i.json);
let last = $input.all(1).map(i => i.json);
if (cur.length === 0 || last.length === 0) {
  const all = $input.all();
  const c = []; const l = [];
  for (const it of all) {
    const idx = it.pairedItem && typeof it.pairedItem.input === 'number' ? it.pairedItem.input : null;
    if (idx === 0) c.push(it.json); else if (idx === 1) l.push(it.json);
  }
  if (cur.length === 0) cur = c; if (last.length === 0) last = l;
}
function norm(v){return (v??'').toString().trim();}
function num(v){ if(v===null||v===undefined||v==='') return 0; let s=String(v).trim(); s=s.replace(/[\s$€¥,]/g,''); const n=Number(s); return isNaN(n)?0:n; }
function getCostUSD(row){ return num(row.CostUSD); }
function key(row){ return norm(row.Resource); }
function aggregate(arr){ const m = new Map(); for (const r of arr){ const k = key(r); if (!k) continue; const c = getCostUSD(r); const prev = m.get(k); if (prev){ prev.cost += c; prev.row = r; } else { m.set(k, { cost: c, row: r }); } } return m; }
const curMap = aggregate(cur);
const lastMap = aggregate(last);
const keys = new Set([...curMap.keys(), ...lastMap.keys()]);
const out = [];
for (const k of keys){
  const c = curMap.get(k); const l = lastMap.get(k);
  const base = (c ? c.row : l.row);
  const curCost = c ? c.cost : 0;
  const lastCost = l ? l.cost : 0;
  const delta = +(curCost - lastCost).toFixed(6);
  let remark;
  if (!l) remark = 'NEW'; else if (!c) remark = 'REMOVED'; else if (delta > 0) remark = 'UP'; else if (delta < 0) remark = 'DOWN'; else remark = 'SAME';
  out.push({
    Resource: k,
    SubscriptionName: norm(base.SubscriptionName || ''),
    ResourceGroup: norm(base.ResourceGroup || ''),
    CurrentCostUSD: +curCost.toFixed(6),
    LastCostUSD: +lastCost.toFixed(6),
    DeltaUSD: delta,
    Remark: remark
  });
}
out.sort((a,b)=> Math.abs(b.DeltaUSD) - Math.abs(a.DeltaUSD));
return out.map(x=> ({ json: x }));
```

## 方案 B：低代码管线（可选）
- 两个 `Item Lists` 节点对各自输入执行 `Aggregate`：
  - 分组键：`Resource`
  - 聚合：`sum(CostUSD)`，保留 `SubscriptionName/ResourceGroup` 的第一次值
- 一个 `Code` 节点执行“按 Resource 关联”和“增删标记”，只负责二路合并与差值计算；相比方案 A，代码更短，但依赖 `Item Lists` 的配置正确。

## 验证用例
- 样例：`vm104-prod-csc-retail-web03`
  - 10 月 `CostUSD=18.93480771`
  - 11 月 `CostUSD=18.81247253`
  - 预期：`DeltaUSD≈-0.12233518`、`Remark=DOWN`
- 执行后在 CSV 中检索该资源，应与上述一致；新增/移除将分别显示 `NEW/REMOVED`。

## 执行步骤
1. 在 n8n 打开该工作流，确保“读取本月账单”连接到“对比分析”的输入 0，“读取上月账单”连接到输入 1。
2. 将“对比分析”节点代码替换为方案 A 的版本并保存。
3. 使用 `current_month=11月.csv`、`last_month=10月.csv` 执行测试并检查生成 CSV。