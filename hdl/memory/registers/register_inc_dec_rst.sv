module register_inc_dec_rst #(
    parameter N = 4
) (
    input  logic         inc, dec, rst,
    output logic [N-1:0] data,
);

    always_ff @(posedge inc, posedge dec, posedge rst) begin
        if (rst) data <= 0;
        else if (inc) data <= data + 1;
        else if (dec) data <= data - 1;
    end

endmodule