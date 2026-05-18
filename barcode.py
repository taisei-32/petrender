import os
os.environ["PYOPENGL_PLATFORM"] = "osmesa"
import math
import numpy as np
from PIL import Image
import pyrender
import trimesh
import subprocess
import os

def calc_check_digit(code12):
    total = sum(int(d) * (1 if i % 2 == 0 else 3) for i, d in enumerate(code12))
    return (10 - (total % 10)) % 10

def gen_barcode(data):
    check_digit = calc_check_digit(data)
    data = data + str(check_digit)
    height = 30
    scale = 4
    os.makedirs("gen_imgs", exist_ok=True)
    img_path = f"gen_imgs/{data}.png"

    try:
        result = subprocess.run([
            "zint",
            "--barcode=13",
            f"--data={data}",
            f"--height={height}",
            f"--scale={scale}",
            f"--output={img_path}"
        ], capture_output=True, text=True)
        print(data)
        if result.stderr:
            print(f"{data} エラー: {result.stderr}")
    except Exception as e:
        print(f"実行エラー: {e}")
    return img_path, data

def make_barcode_patch(radius, height, barcode_angle_deg, segments=128):
        half_angle = math.radians(barcode_angle_deg / 2.0)

        vertices = []
        faces = []
        uvs = []

        for i in range(segments + 1):
            theta = -half_angle + (2 * half_angle * i / segments)


            r = radius

            x = r * math.cos(theta)
            y = r * math.sin(theta)

            u = i / segments

            vertices.append([x, y, -height / 2])
            vertices.append([x, y,  height / 2])

            uvs.append([u, 0.0])
            uvs.append([u, 1.0])

        for i in range(segments):
            b = i * 2
            faces.append([b, b + 2, b + 3])
            faces.append([b, b + 3, b + 1])

        return (
            np.array(vertices, dtype=np.float32),
            np.array(faces, dtype=np.int32),
            np.array(uvs, dtype=np.float32)
        )

def render_scene(barcode_angle_deg, camera_offset_deg, barcode_img,
                 CYLINDER_RADIUS, CYLINDER_HEIGHT,
                 CAMERA_DISTANCE, IMG_WIDTH, IMG_HEIGHT, camera_fov):

    scene = pyrender.Scene(
        bg_color=[255, 255, 255, 255],
        ambient_light=[0.8, 0.8, 0.8]
    )

    cylinder = trimesh.creation.cylinder(
        radius=CYLINDER_RADIUS,
        height=CYLINDER_HEIGHT,
        sections=128
    )
    white_material = pyrender.MetallicRoughnessMaterial(
        baseColorFactor=[1.0, 1.0, 1.0, 1.0],
        roughnessFactor=1.0,
        metallicFactor=0.0
    )
    scene.add(pyrender.Mesh.from_trimesh(cylinder, material=white_material, smooth=False))

    vertices, faces, uvs = make_barcode_patch(
        CYLINDER_RADIUS, CYLINDER_HEIGHT, barcode_angle_deg
    )
    patch = trimesh.Trimesh(vertices=vertices, faces=faces, process=False)
    patch.visual = trimesh.visual.TextureVisuals(uv=uvs, image=barcode_img)
    scene.add(pyrender.Mesh.from_trimesh(patch, smooth=False))

    offset_rad = math.radians(camera_offset_deg)
    
    cam_x = CAMERA_DISTANCE * math.cos(offset_rad)
    cam_y = CAMERA_DISTANCE * math.sin(offset_rad)

    cam_pose = np.array([
        [-math.sin(offset_rad),  0,  math.cos(offset_rad),  cam_x],
        [ math.cos(offset_rad),  0,  math.sin(offset_rad),  cam_y],
        [ 0,                     1,  0,                      0    ],
        [ 0,                     0,  0,                      1    ],
    ], dtype=np.float64)

    camera = pyrender.PerspectiveCamera(yfov=math.radians(camera_fov))
    scene.add(camera, pose=cam_pose)

    light = pyrender.DirectionalLight(color=[1.0, 1.0, 1.0], intensity=3.0)
    scene.add(light, pose=cam_pose)

    renderer = pyrender.OffscreenRenderer(IMG_WIDTH, IMG_HEIGHT)
    color, _ = renderer.render(scene, flags=pyrender.RenderFlags.RGBA)
    renderer.delete()

    return np.array(color, dtype=np.uint8)


def render(img_path, data, camera_fov):
    OUTPUT_DIR  = f"render_imgs/{data}"
    ANGLE_START = 10                 
    ANGLE_END   = 80                
    ANGLE_STEP  = 1   

    OFFSET_START = -30
    OFFSET_END   = 30
    OFFSET_STEP  = 1

    CYLINDER_RADIUS  = 1.0
    CYLINDER_HEIGHT  = 0.5
    CAMERA_DISTANCE  = 3.0
    IMG_WIDTH        = 1280
    IMG_HEIGHT       = 1060

    os.makedirs(OUTPUT_DIR, exist_ok=True)

    barcode_img = Image.open(img_path).convert("RGBA")
    angles = list(range(ANGLE_START, ANGLE_END + 1, ANGLE_STEP))

    for angle in angles:
        for offset in range(OFFSET_START, OFFSET_END + 1, OFFSET_STEP):
            img_array = render_scene(
                angle, offset, barcode_img,
                CYLINDER_RADIUS, CYLINDER_HEIGHT,
                CAMERA_DISTANCE, IMG_WIDTH, IMG_HEIGHT, camera_fov
            )
            img = Image.fromarray(img_array).rotate(0, expand=True)
            out_path = os.path.join(OUTPUT_DIR, f"barcode_{angle}_{offset:03d}deg.png")
            img.save(out_path)
            print(f"{out_path}")

def main():
    raw_data = 495912740082
    offset = 23
    # 12桁
    subprocess.run(["sed", "-i", r"s/\r//", "analyze.sh"])
    for data in range(raw_data, raw_data+offset):
        data = str(data)
        img_path, data = gen_barcode(data)
        render(img_path, data, 43)
    for data in range(raw_data, raw_data+offset):
        data = str(data)
        check_digit = calc_check_digit(data)
        data = data + str(check_digit)  
        subprocess.run(["bash", "analyze.sh", data.strip()], check=True)
    subprocess.run(["./classify_log","analyze/result_log/filter/zbar","analyze/result_log/classify/zbar.txt"], check=True) 
    subprocess.run(["./classify_log","analyze/result_log/filter/zxing","analyze/result_log/classify/zxing.txt"], check=True) 
    subprocess.run(["./summary_log","analyze/result_log/reg/zbar","analyze/result_log/zbar.txt"], check=True)
    subprocess.run(["./summary_log","analyze/result_log/reg/zxing","analyze/result_log/zxing.txt"], check=True)
if __name__ == "__main__":
    main()