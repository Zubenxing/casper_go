import csv
from collections import defaultdict
from pathlib import Path

BASE = Path(__file__).resolve().parent

last_path = BASE / "10月.csv"  # 上个月
cur_path = BASE / "11月.csv"   # 本月


def load_agg(path: Path):
    agg = defaultdict(float)
    with path.open("r", encoding="utf-8-sig", newline="") as f:
        reader = csv.DictReader(f)
        for row in reader:
            resource = (row.get("Resource") or "").strip()
            sub = (row.get("SubscriptionName") or "").strip()
            rg = (row.get("ResourceGroup") or "").strip()
            if not resource and not sub and not rg:
                continue
            cost_str = (row.get("CostUSD") or "").strip()
            if not cost_str:
                value = 0.0
            else:
                # 去掉可能的逗号和空格
                value = float(cost_str.replace(",", ""))
            key = (resource, sub, rg)
            agg[key] += value
    return agg


def main():
    last = load_agg(last_path)
    cur = load_agg(cur_path)

    keys = set(last.keys()) | set(cur.keys())
    rows = []
    for k in keys:
        resource, sub, rg = k
        last_v = last.get(k, 0.0)
        cur_v = cur.get(k, 0.0)
        delta = cur_v - last_v
        if last_v == 0 and cur_v > 0:
            remark = "NEW"
        elif cur_v == 0 and last_v > 0:
            remark = "REMOVED"
        elif delta > 0:
            remark = "UP"
        elif delta < 0:
            remark = "DOWN"
        else:
            remark = "SAME"
        rows.append(
            {
                "Resource": resource,
                "SubscriptionName": sub,
                "ResourceGroup": rg,
                "CurrentCostUSD": round(cur_v, 6),
                "LastCostUSD": round(last_v, 6),
                "DeltaUSD": round(delta, 6),
                "Remark": remark,
            }
        )

    # 按变动绝对值排序，方便人工检查
    rows.sort(key=lambda r: (-abs(r["DeltaUSD"]), r["Resource"]))

    out_path = BASE / "azure-bill-comparison-python.csv"
    with out_path.open("w", encoding="utf-8", newline="") as f:
        writer = csv.DictWriter(
            f,
            fieldnames=[
                "Resource",
                "SubscriptionName",
                "ResourceGroup",
                "CurrentCostUSD",
                "LastCostUSD",
                "DeltaUSD",
                "Remark",
            ],
        )
        writer.writeheader()
        writer.writerows(rows)

    print(f"写出对比结果: {out_path}")
    # 简单打印几个关注的资源
    watch = {
        "vm127-prod-tcn",
        "prod-sage-db01",
        "seatunnel-cluster-test1",
    }
    for r in rows:
        if r["Resource"] in watch:
            print(r)


if __name__ == "__main__":
    main()
